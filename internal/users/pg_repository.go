package users

import (
	"context"
	"time"

	"github.com/hiennguyen9874/go-boilerplate/internal"
	"github.com/hiennguyen9874/go-boilerplate/internal/models"
)

type UserPgRepository interface {
	internal.PgRepository[models.User]
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	UpdatePassword(ctx context.Context, exp *models.User, newPassword string) (*models.User, error)
	UpdateVerificationCode(
		ctx context.Context,
		exp *models.User,
		newVerificationCode string,
	) (*models.User, error)
	UpdateVerification(
		ctx context.Context,
		exp *models.User,
		newVerificationCode string,
		newVerified bool,
	) (*models.User, error)
	GetByVerificationCode(ctx context.Context, verificationCode string) (*models.User, error)
	UpdatePasswordReset(
		ctx context.Context,
		exp *models.User,
		passwordResetToken string,
		passwordResetAt time.Time,
	) (*models.User, error)
	GetByResetToken(ctx context.Context, resetToken string) (*models.User, error)
	GetByResetTokenResetAt(ctx context.Context, resetToken string, resetAt time.Time) (*models.User, error)
	UpdatePasswordResetToken(
		ctx context.Context,
		exp *models.User,
		newPassword string,
		resetToken string,
	) (*models.User, error)
}
