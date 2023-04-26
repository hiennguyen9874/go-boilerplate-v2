package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/hiennguyen9874/go-boilerplate-v2/config"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/models"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/usecase"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/users"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/worker"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/cryptpass"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/emailTemplates"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/httpErrors"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/jwt"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/logger"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/secureRandom"
)

type userUseCase struct {
	usecase.UseCase[models.User]
	pgRepo                 users.UserPgRepository
	redisRepo              users.UserRedisRepository
	emailTemplateGenerator emailTemplates.EmailTemplatesGenerator
	redisTaskDistributor   users.UserRedisTaskDistributor
}

func CreateUserUseCaseI(
	pgRepo users.UserPgRepository,
	redisRepo users.UserRedisRepository,
	redisTaskDistributor users.UserRedisTaskDistributor,
	cfg *config.Config,
	logger logger.Logger,
) users.UserUseCaseI {
	return &userUseCase{
		UseCase:                usecase.CreateUseCase[models.User](pgRepo, cfg, logger),
		pgRepo:                 pgRepo,
		redisRepo:              redisRepo,
		emailTemplateGenerator: emailTemplates.NewEmailTemplatesGenerator(cfg),
		redisTaskDistributor:   redisTaskDistributor,
	}
}

func (u *userUseCase) Get(ctx context.Context, id uuid.UUID) (*models.User, error) {
	cachedUser, err := u.redisRepo.Get(ctx, u.GenerateRedisUserKey(id))
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

	if err = u.redisRepo.Create(ctx, u.GenerateRedisUserKey(id), user, 3600); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUseCase) Delete(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := u.pgRepo.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.Delete(ctx, u.GenerateRedisUserKey(id)); err != nil {
		return nil, err
	}

	if err = u.redisRepo.Delete(ctx, u.GenerateRedisRefreshTokenKey(id)); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUseCase) Update(
	ctx context.Context,
	id uuid.UUID,
	values map[string]interface{},
) (*models.User, error) {
	obj, err := u.Get(ctx, id)
	if err != nil || obj == nil {
		return nil, err
	}

	user, err := u.pgRepo.Update(ctx, obj, values)
	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.Delete(ctx, u.GenerateRedisUserKey(id)); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUseCase) Create(ctx context.Context, exp *models.User) (*models.User, error) {
	exp.Email = strings.ToLower(strings.TrimSpace(exp.Email))
	exp.Password = strings.TrimSpace(exp.Password)

	hashedPassword, err := cryptpass.HashPassword(exp.Password)
	if err != nil {
		return nil, err
	}
	exp.Password = hashedPassword

	user, err := u.pgRepo.Create(ctx, exp)
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
	updatedUser, err := u.pgRepo.UpdateVerificationCode(ctx, user, verificationCode)
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
		From:      u.Cfg.Email.From,
		To:        updatedUser.Email,
		Subject:   u.Cfg.Email.VerificationSubject,
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

func (u *userUseCase) CreateUser(ctx context.Context, exp *models.User, confirmPassword string) (*models.User, error) {
	if exp.Password != confirmPassword {
		return nil, httpErrors.ErrValidation(errors.New("password do not match"))
	}
	return u.Create(ctx, exp)
}

func (u *userUseCase) createToken(ctx context.Context, exp models.User) (string, string, error) {
	accessToken, err := jwt.CreateAccessTokenRS256(
		exp.Id.String(),
		exp.Email,
		u.Cfg.Jwt.AccessTokenPrivateKey,
		u.Cfg.Jwt.AccessTokenExpireDuration*int64(time.Minute),
		u.Cfg.Jwt.Issuer,
	)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.CreateAccessTokenRS256(
		exp.Id.String(),
		exp.Email,
		u.Cfg.Jwt.RefreshTokenPrivateKey,
		u.Cfg.Jwt.RefreshTokenExpireDuration*int64(time.Minute),
		u.Cfg.Jwt.Issuer,
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

func (u *userUseCase) IsActive(ctx context.Context, exp models.User) bool {
	return exp.IsActive
}

func (u *userUseCase) IsSuper(ctx context.Context, exp models.User) bool {
	return exp.IsSuperUser
}

func (u *userUseCase) CreateSuperUserIfNotExist(ctx context.Context) (bool, error) {
	user, err := u.pgRepo.GetByEmail(ctx, u.Cfg.FirstSuperUser.Email)

	if err != nil || user == nil {
		_, err := u.Create(ctx, &models.User{
			Name:        u.Cfg.FirstSuperUser.Name,
			Email:       u.Cfg.FirstSuperUser.Email,
			Password:    u.Cfg.FirstSuperUser.Password,
			IsActive:    true,
			IsSuperUser: true,
			Verified:    true,
		})
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func (u *userUseCase) UpdatePassword(
	ctx context.Context,
	id uuid.UUID,
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

	updatedUser, err := u.pgRepo.UpdatePassword(ctx, user, hashedPassword)
	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.Delete(ctx, u.GenerateRedisUserKey(id)); err != nil {
		return nil, err
	}

	if err = u.redisRepo.Delete(ctx, u.GenerateRedisRefreshTokenKey(id)); err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (u *userUseCase) ParseIdFromRefreshToken(
	ctx context.Context,
	refreshToken string,
) (uuid.UUID, error) {
	id, _, err := jwt.ParseTokenRS256(refreshToken, u.Cfg.Jwt.RefreshTokenPublicKey)
	if err != nil {
		return uuid.UUID{}, err
	}

	idParsed, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{},
			httpErrors.ErrInvalidJWTClaims(errors.New("can not convert id to uuid from id in token"))
	}

	return idParsed, nil
}

func (u *userUseCase) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	idParsed, err := u.ParseIdFromRefreshToken(ctx, refreshToken)
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
	idParsed, err := u.ParseIdFromRefreshToken(ctx, refreshToken)
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

func (u *userUseCase) LogoutAll(ctx context.Context, id uuid.UUID) error {
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

	updatedUser, err := u.pgRepo.UpdateVerification(ctx, user, "", true)
	if err != nil {
		return err
	}

	if err = u.redisRepo.Delete(ctx, u.GenerateRedisUserKey(updatedUser.Id)); err != nil {
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
		user,
		resetToken,
		time.Now().Add(time.Minute*15),
	)
	if err != nil {
		return err
	}
	if err = u.redisRepo.Delete(ctx, u.GenerateRedisUserKey(updatedUser.Id)); err != nil {
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
		From:      u.Cfg.Email.From,
		To:        updatedUser.Email,
		Subject:   u.Cfg.Email.ResetSubject,
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
		user,
		hashedPassword,
		"",
	)
	if err != nil {
		return err
	}

	if err = u.redisRepo.Delete(ctx, u.GenerateRedisUserKey(updatedUser.Id)); err != nil {
		return err
	}

	if err = u.redisRepo.Delete(ctx, u.GenerateRedisRefreshTokenKey(updatedUser.Id)); err != nil {
		return err
	}

	return nil
}

func (u *userUseCase) GenerateRedisUserKey(id uuid.UUID) string {
	return fmt.Sprintf("%s:%s", models.User{}.TableName(), id.String())
}

func (u *userUseCase) GenerateRedisRefreshTokenKey(id uuid.UUID) string {
	return fmt.Sprintf("RefreshToken:%s", id.String())
}
