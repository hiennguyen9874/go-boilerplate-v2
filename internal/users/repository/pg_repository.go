package repository

import (
	"context"
	"time"

	"github.com/hiennguyen9874/go-boilerplate/internal/models"
	"github.com/hiennguyen9874/go-boilerplate/internal/repository"
	"github.com/hiennguyen9874/go-boilerplate/internal/users"
	"gorm.io/gorm"
)

type UserPgRepo struct {
	repository.PgRepo[models.User]
}

func CreateUserPgRepository(db *gorm.DB) users.UserPgRepository {
	return &UserPgRepo{
		PgRepo: repository.CreatePgRepo[models.User](db),
	}
}

func (r *UserPgRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var obj *models.User
	if result := r.DB.WithContext(ctx).First(&obj, "email = ?", email); result.Error != nil {
		return nil, result.Error
	}
	return obj, nil
}

func (r *UserPgRepo) UpdatePassword(
	ctx context.Context,
	exp *models.User,
	newPassword string,
) (*models.User, error) {
	if result := r.DB.WithContext(ctx).Model(&exp).Select("password").
		Updates(map[string]interface{}{"password": newPassword}); result.Error != nil {
		return nil, result.Error
	}
	return exp, nil
}

func (r *UserPgRepo) UpdateVerificationCode(
	ctx context.Context,
	exp *models.User,
	newVerificationCode string,
) (*models.User, error) {
	if result := r.DB.WithContext(ctx).Model(&exp).Select("verification_code").
		Updates(map[string]interface{}{"verification_code": newVerificationCode}); result.Error != nil {
		return nil, result.Error
	}
	return exp, nil
}

func (r *UserPgRepo) UpdateVerification(
	ctx context.Context,
	exp *models.User,
	newVerificationCode string,
	newVerified bool,
) (*models.User, error) {
	if result := r.DB.WithContext(ctx).Model(&exp).Select("verification_code", "verified").
		Updates(map[string]interface{}{
			"verification_code": newVerificationCode,
			"verified":          newVerified,
		}); result.Error != nil {
		return nil, result.Error
	}
	return exp, nil
}

func (r *UserPgRepo) GetByVerificationCode(
	ctx context.Context,
	verificationCode string,
) (*models.User, error) {
	var obj *models.User
	if result := r.DB.WithContext(ctx).First(&obj, "verification_code = ?", verificationCode); result.Error != nil {
		return nil, result.Error
	}
	return obj, nil
}

func (r *UserPgRepo) UpdatePasswordReset(
	ctx context.Context,
	exp *models.User,
	passwordResetToken string,
	passwordResetAt time.Time,
) (*models.User, error) {
	if result := r.DB.WithContext(ctx).Model(&exp).Select("password_reset_token", "password_reset_at").
		Updates(map[string]interface{}{
			"password_reset_token": passwordResetToken,
			"password_reset_at":    passwordResetAt,
		}); result.Error != nil {
		return nil, result.Error
	}
	return exp, nil
}

func (r *UserPgRepo) GetByResetToken(ctx context.Context, resetToken string) (*models.User, error) {
	var obj *models.User
	if result := r.DB.WithContext(ctx).First(&obj, "reset_token = ?", resetToken); result.Error != nil {
		return nil, result.Error
	}
	return obj, nil
}

func (r *UserPgRepo) GetByResetTokenResetAt(
	ctx context.Context,
	resetToken string,
	resetAt time.Time,
) (*models.User, error) {
	var obj *models.User
	if result := r.DB.WithContext(ctx).First(&obj, "password_reset_token = ? AND password_reset_at > ?", resetToken, resetAt); result.Error != nil {
		return nil, result.Error
	}
	return obj, nil
}

func (r *UserPgRepo) UpdatePasswordResetToken(
	ctx context.Context,
	exp *models.User,
	newPassword string,
	resetToken string,
) (*models.User, error) {
	if result := r.DB.WithContext(ctx).Model(&exp).Select("password", "password_reset_token").
		Updates(map[string]interface{}{"password": newPassword, "password_reset_token": resetToken}); result.Error != nil {
		return nil, result.Error
	}
	return exp, nil
}
