package users

import (
	"context"

	"github.com/hiennguyen9874/go-boilerplate-v2/internal/models"
)

type UserRedisRepository interface {
	Create(ctx context.Context, key string, obj_in *models.User, seconds int) error
	Get(ctx context.Context, key string) (*models.User, error)
	Delete(ctx context.Context, key string) error
	Sadd(ctx context.Context, key string, value string) error
	Sadds(ctx context.Context, key string, values []string) error
	Srem(ctx context.Context, key string, value string) error
	SIsMember(ctx context.Context, key string, value string) (bool, error)
}
