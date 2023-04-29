package users

import (
	"context"
	"time"

	"github.com/hiennguyen9874/go-boilerplate-v2/internal/models"
)

type UserPgRepository interface {
	Get(ctx context.Context, id uint) (*models.User, error)
	GetMulti(ctx context.Context, offset, limit int) ([]*models.User, error)
	Create(ctx context.Context, obj_in *models.UserCreate) (*models.User, error)
	Delete(ctx context.Context, id uint) (*models.User, error)
	Update(ctx context.Context, id uint, obj_update *models.UserUpdate) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	UpdatePassword(ctx context.Context, id uint, newPassword string) (*models.User, error)
	UpdateVerificationCode(ctx context.Context, id uint, newVerificationCode string) (*models.User, error)
	UpdateVerification(ctx context.Context, id uint, newVerificationCode string, newVerified bool) (*models.User, error)
	GetByVerificationCode(ctx context.Context, verificationCode string) (*models.User, error)
	UpdatePasswordReset(ctx context.Context, id uint, passwordResetToken string, passwordResetAt time.Time) (*models.User, error)
	GetByResetToken(ctx context.Context, resetToken string) (*models.User, error)
	GetByResetTokenResetAt(ctx context.Context, resetToken string, resetAt time.Time) (*models.User, error)
	UpdatePasswordResetToken(ctx context.Context, id uint, newPassword string, resetToken string) (*models.User, error)
}
