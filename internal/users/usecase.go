package users

import (
	"context"

	"github.com/hiennguyen9874/go-boilerplate-v2/internal/models"
)

type UserUseCase interface {
	Create(ctx context.Context, obj_create *models.UserCreate, confirmPassword string) (*models.User, error)
	Get(ctx context.Context, id uint) (*models.User, error)
	GetMulti(ctx context.Context, offset, limit int) ([]*models.User, error)
	Delete(ctx context.Context, id uint) (*models.User, error)
	Update(ctx context.Context, id uint, obj_update *models.UserUpdate) (*models.User, error)
	SignIn(ctx context.Context, email string, password string) (string, string, error)
	CreateSuperUserIfNotExist(context.Context) (bool, error)
	UpdatePassword(ctx context.Context, id uint, oldPassword string, newPassword string, confirmPassword string) (*models.User, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, error)
	Logout(ctx context.Context, refreshToken string) error
	LogoutAllWithId(ctx context.Context, id uint) error
	LogoutAllWithToken(ctx context.Context, refreshToken string) error
	Verify(ctx context.Context, verificationCode string) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, resetToken string, newPassword string, confirmPassword string) error
}
