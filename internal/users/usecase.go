package users

import (
	"context"

	"github.com/google/uuid"
	"github.com/hiennguyen9874/go-boilerplate/internal"
	"github.com/hiennguyen9874/go-boilerplate/internal/models"
)

type UserUseCaseI interface {
	internal.UseCaseI[models.User]
	CreateUser(ctx context.Context, exp *models.User, confirmPassword string) (*models.User, error)
	SignIn(ctx context.Context, email string, password string) (string, string, error)
	IsActive(ctx context.Context, exp models.User) bool
	IsSuper(ctx context.Context, exp models.User) bool
	CreateSuperUserIfNotExist(context.Context) (bool, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, oldPassword string, newPassword string, confirmPassword string) (*models.User, error)
	ParseIdFromRefreshToken(ctx context.Context, refreshToken string) (uuid.UUID, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, error)
	GenerateRedisUserKey(id uuid.UUID) string
	GenerateRedisRefreshTokenKey(id uuid.UUID) string
	Logout(ctx context.Context, refreshToken string) error
	LogoutAll(ctx context.Context, id uuid.UUID) error
	Verify(ctx context.Context, verificationCode string) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, resetToken string, newPassword string, confirmPassword string) error
}
