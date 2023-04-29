package usecase

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hibiken/asynq"
	"github.com/hiennguyen9874/go-boilerplate-v2/config"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/models"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/users"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/worker"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/cryptpass"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/emailTemplates"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/httpErrors"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/jwt"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/logger"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/newPointer"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/secureRandom"
)

type userUseCase struct {
	pgRepo                 users.UserPgRepository
	redisRepo              users.UserRedisRepository
	emailTemplateGenerator emailTemplates.EmailTemplatesGenerator
	redisTaskDistributor   users.UserRedisTaskDistributor
	cfg                    *config.Config
	logger                 logger.Logger
}

func CreateUserUseCase(
	pgRepo users.UserPgRepository,
	redisRepo users.UserRedisRepository,
	redisTaskDistributor users.UserRedisTaskDistributor,
	cfg *config.Config,
	logger logger.Logger,
) users.UserUseCase {
	return &userUseCase{
		pgRepo:                 pgRepo,
		redisRepo:              redisRepo,
		emailTemplateGenerator: emailTemplates.NewEmailTemplatesGenerator(cfg),
		redisTaskDistributor:   redisTaskDistributor,
		cfg:                    cfg,
		logger:                 logger,
	}
}

func (u *userUseCase) Create(ctx context.Context, obj_create *models.UserCreate, confirmPassword string) (*models.User, error) {
	if obj_create.Password != confirmPassword {
		return nil, httpErrors.ErrValidation(errors.New("password do not match"))
	}

	obj_create.Email = strings.ToLower(strings.TrimSpace(obj_create.Email))
	obj_create.Password = strings.TrimSpace(obj_create.Password)

	hashedPassword, err := cryptpass.HashPassword(obj_create.Password)
	if err != nil {
		return nil, err
	}
	obj_create.Password = hashedPassword

	user, err := u.pgRepo.Create(ctx, obj_create)
	if err != nil {
		return nil, err
	}

	if user.Verified {
		return user, nil
	}

	verificationCode, err := secureRandom.RandomHex(16)
	if err != nil {
		return nil, err
	}

	// Update user in database
	updatedUser, err := u.pgRepo.UpdateVerificationCode(ctx, user.Id, verificationCode)
	if err != nil {
		return nil, err
	}

	bodyHtml, bodyPlain, err := u.emailTemplateGenerator.GenerateVerificationCodeTemplate(
		ctx,
		updatedUser.Name,
		fmt.Sprintf(
			"http://localhost:5000/auth/verifyemail?code=%s",
			verificationCode,
		),
	)
	if err != nil {
		return nil, err
	}

	err = u.redisTaskDistributor.DistributeTaskSendEmail(ctx, &users.PayloadSendEmail{
		From:      u.cfg.Email.From,
		To:        updatedUser.Email,
		Subject:   u.cfg.Email.VerificationSubject,
		BodyHtml:  bodyHtml,
		BodyPlain: bodyPlain,
	}, []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}...)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (u *userUseCase) Get(ctx context.Context, id uint) (*models.User, error) {
	cachedUser, err := u.redisRepo.Get(ctx, u.generateRedisUserKey(id))
	if err != nil {
		return nil, err
	}

	if cachedUser != nil {
		return cachedUser, nil
	}

	user, err := u.pgRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.Create(ctx, u.generateRedisUserKey(id), user, 3600); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUseCase) GetMulti(ctx context.Context, offset, limit int) ([]*models.User, error) {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 50
	}
	return u.pgRepo.GetMulti(ctx, offset, limit)
}

func (u *userUseCase) Delete(ctx context.Context, id uint) (*models.User, error) {
	user, err := u.pgRepo.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.Delete(ctx, u.generateRedisUserKey(id)); err != nil {
		return nil, err
	}

	if err = u.redisRepo.Delete(ctx, u.GenerateRedisRefreshTokenKey(id)); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUseCase) Update(
	ctx context.Context,
	id uint,
	obj_update *models.UserUpdate,
) (*models.User, error) {
	obj, err := u.Get(ctx, id)
	if err != nil || obj == nil {
		return nil, err
	}

	user, err := u.pgRepo.Update(ctx, obj.Id, obj_update)
	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.Delete(ctx, u.generateRedisUserKey(id)); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUseCase) createToken(ctx context.Context, obj models.User) (string, string, error) {
	accessToken, err := jwt.CreateAccessTokenRS256(
		strconv.FormatUint(uint64(obj.Id), 10),
		obj.Email,
		u.cfg.Jwt.AccessTokenPrivateKey,
		u.cfg.Jwt.AccessTokenExpireDuration*int64(time.Minute),
		u.cfg.Jwt.Issuer,
	)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.CreateAccessTokenRS256(
		strconv.FormatUint(uint64(obj.Id), 10),
		obj.Email,
		u.cfg.Jwt.RefreshTokenPrivateKey,
		u.cfg.Jwt.RefreshTokenExpireDuration*int64(time.Minute),
		u.cfg.Jwt.Issuer,
	)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (u *userUseCase) SignIn(ctx context.Context, email string, password string) (string, string, error) {
	user, err := u.pgRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", "", httpErrors.ErrNotFound(err)
	}

	if !cryptpass.ComparePassword(password, user.Password) {
		return "", "", httpErrors.ErrWrongPassword(errors.New("wrong password"))
	}

	accessToken, refreshToken, err := u.createToken(ctx, *user)
	if err != nil {
		return "", "", err
	}

	if err = u.redisRepo.Sadd(
		ctx,
		u.GenerateRedisRefreshTokenKey(user.Id),
		refreshToken,
	); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (u *userUseCase) CreateSuperUserIfNotExist(ctx context.Context) (bool, error) {
	user, err := u.pgRepo.GetByEmail(ctx, u.cfg.FirstSuperUser.Email)

	if err != nil || user == nil {
		_, err := u.Create(ctx, &models.UserCreate{
			Name:        u.cfg.FirstSuperUser.Name,
			Email:       u.cfg.FirstSuperUser.Email,
			Password:    u.cfg.FirstSuperUser.Password,
			IsActive:    newPointer.NewBoolean(true),
			IsSuperUser: newPointer.NewBoolean(true),
		}, u.cfg.FirstSuperUser.Password)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func (u *userUseCase) UpdatePassword(
	ctx context.Context,
	id uint,
	oldPassword string,
	newPassword string,
	confirmPassword string,
) (*models.User, error) {
	if newPassword != confirmPassword {
		return nil, httpErrors.ErrValidation(errors.New("password do not match"))
	}

	user, err := u.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if !cryptpass.ComparePassword(oldPassword, user.Password) {
		return nil, httpErrors.ErrWrongPassword(errors.New("old password and new password not same"))
	}

	hashedPassword, err := cryptpass.HashPassword(newPassword)
	if err != nil {
		return nil, err
	}

	updatedUser, err := u.pgRepo.UpdatePassword(ctx, user.Id, hashedPassword)
	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.Delete(ctx, u.generateRedisUserKey(id)); err != nil {
		return nil, err
	}

	if err = u.redisRepo.Delete(ctx, u.GenerateRedisRefreshTokenKey(id)); err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (u *userUseCase) parseIdFromRefreshToken(
	ctx context.Context,
	refreshToken string,
) (uint, error) {
	id, _, err := jwt.ParseTokenRS256(refreshToken, u.cfg.Jwt.RefreshTokenPublicKey)
	if err != nil {
		return 0, err
	}

	idParsed, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return 0,
			httpErrors.ErrInvalidJWTClaims(errors.New("can not convert id to uuid from id in token"))
	}

	return uint(idParsed), nil
}

func (u *userUseCase) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	idParsed, err := u.parseIdFromRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	isMember, err := u.redisRepo.SIsMember(
		ctx,
		u.GenerateRedisRefreshTokenKey(idParsed),
		refreshToken,
	)
	if err != nil {
		return "", "", err
	}

	if !isMember {
		return "", "",
			httpErrors.ErrNotFoundRefreshTokenRedis(errors.New("not found refresh token in redis"))
	}

	if err = u.redisRepo.Srem(
		ctx,
		u.GenerateRedisRefreshTokenKey(idParsed),
		refreshToken,
	); err != nil {
		return "", "", err
	}

	user, err := u.Get(ctx, idParsed)
	if err != nil {
		return "", "", err
	}

	accessToken, refreshToken, err := u.createToken(ctx, *user)
	if err != nil {
		return "", "", err
	}

	if err = u.redisRepo.Sadd(
		ctx,
		u.GenerateRedisRefreshTokenKey(user.Id),
		refreshToken,
	); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, err
}

func (u *userUseCase) Logout(ctx context.Context, refreshToken string) error {
	idParsed, err := u.parseIdFromRefreshToken(ctx, refreshToken)
	if err != nil {
		return err
	}

	if err = u.redisRepo.Srem(
		ctx,
		u.GenerateRedisRefreshTokenKey(idParsed),
		refreshToken,
	); err != nil {
		return err
	}

	return nil
}

func (u *userUseCase) LogoutAllWithId(ctx context.Context, id uint) error {
	if err := u.redisRepo.Delete(ctx, u.GenerateRedisRefreshTokenKey(id)); err != nil {
		return err
	}

	return nil
}

func (u *userUseCase) LogoutAllWithToken(ctx context.Context, refreshToken string) error {
	id, err := u.parseIdFromRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil
	}

	if err := u.redisRepo.Delete(ctx, u.GenerateRedisRefreshTokenKey(id)); err != nil {
		return err
	}

	return nil
}

func (u *userUseCase) Verify(ctx context.Context, verificationCode string) error {
	user, err := u.pgRepo.GetByVerificationCode(ctx, verificationCode)
	if err != nil {
		return err
	}

	if user.Verified {
		return httpErrors.ErrUserAlreadyVerified(errors.New("user already verified"))
	}

	updatedUser, err := u.pgRepo.UpdateVerification(ctx, user.Id, "", true)
	if err != nil {
		return err
	}

	if err = u.redisRepo.Delete(ctx, u.generateRedisUserKey(updatedUser.Id)); err != nil {
		return err
	}

	return nil
}

func (u *userUseCase) ForgotPassword(ctx context.Context, email string) error {
	user, err := u.pgRepo.GetByEmail(ctx, email)

	if err != nil {
		return httpErrors.ErrNotFound(err)
	}

	if !user.Verified {
		return httpErrors.ErrUserNotVerified(errors.New("user not verified"))
	}

	resetToken, err := secureRandom.RandomHex(16)
	if err != nil {
		return err
	}

	updatedUser, err := u.pgRepo.UpdatePasswordReset(
		ctx,
		user.Id,
		resetToken,
		time.Now().Add(time.Minute*15),
	)
	if err != nil {
		return err
	}
	if err = u.redisRepo.Delete(ctx, u.generateRedisUserKey(updatedUser.Id)); err != nil {
		return err
	}

	bodyHtml, bodyPlain, err := u.emailTemplateGenerator.GeneratePasswordResetTemplate(
		ctx,
		updatedUser.Name,
		fmt.Sprintf("http://localhost:5000/auth/resetpassword?code=%s", resetToken),
	)
	if err != nil {
		return err
	}

	err = u.redisTaskDistributor.DistributeTaskSendEmail(ctx, &users.PayloadSendEmail{
		From:      u.cfg.Email.From,
		To:        updatedUser.Email,
		Subject:   u.cfg.Email.ResetSubject,
		BodyHtml:  bodyHtml,
		BodyPlain: bodyPlain,
	}, []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}...)
	if err != nil {
		return err
	}

	return nil
}

func (u *userUseCase) ResetPassword(
	ctx context.Context,
	resetToken string,
	newPassword string,
	confirmPassword string,
) error {
	if newPassword != confirmPassword {
		return httpErrors.ErrValidation(errors.New("password do not match"))
	}

	user, err := u.pgRepo.GetByResetTokenResetAt(ctx, resetToken, time.Now())
	if err != nil {
		return err
	}

	hashedPassword, err := cryptpass.HashPassword(newPassword)
	if err != nil {
		return err
	}

	updatedUser, err := u.pgRepo.UpdatePasswordResetToken(
		ctx,
		user.Id,
		hashedPassword,
		"",
	)
	if err != nil {
		return err
	}

	if err = u.redisRepo.Delete(ctx, u.generateRedisUserKey(updatedUser.Id)); err != nil {
		return err
	}

	if err = u.redisRepo.Delete(ctx, u.GenerateRedisRefreshTokenKey(updatedUser.Id)); err != nil {
		return err
	}

	return nil
}

func (u *userUseCase) generateRedisUserKey(id uint) string {
	return fmt.Sprintf("Cache:User:%v", id)
}

func (u *userUseCase) GenerateRedisRefreshTokenKey(id uint) string {
	return fmt.Sprintf("RefreshToken:%v", id)
}
