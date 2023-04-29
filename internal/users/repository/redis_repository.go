package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/hiennguyen9874/go-boilerplate-v2/internal/models"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/users"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/httpErrors"
	"github.com/redis/go-redis/v9"
)

type UserRedisRepo struct {
	redisClient *redis.Client
}

func CreateUserRedisRepository(redisClient *redis.Client) users.UserRedisRepository {
	return &UserRedisRepo{redisClient: redisClient}
}

func (r *UserRedisRepo) Create(ctx context.Context, key string, obj_in *models.User, seconds int) error {
	objBytes, err := json.Marshal(obj_in)
	if err != nil {
		return httpErrors.ErrJson(err)
	}

	if err = r.redisClient.Set(ctx, key, objBytes, time.Second*time.Duration(seconds)).Err(); err != nil {
		// TODO: Using httpErrors
		return err
	}
	return nil
}

func (r *UserRedisRepo) Get(ctx context.Context, key string) (*models.User, error) {
	objBytes, err := r.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	var obj models.User

	if err = json.Unmarshal(objBytes, &obj); err != nil {
		return nil, httpErrors.ErrJson(err)
	}

	return &obj, nil
}

func (r *UserRedisRepo) Delete(ctx context.Context, key string) error {
	if err := r.redisClient.Del(ctx, key).Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		// TODO: Using httpErrors
		return err
	}
	return nil
}

func (r *UserRedisRepo) Sadd(ctx context.Context, key string, value string) error {
	if err := r.redisClient.SAdd(ctx, key, value).Err(); err != nil {
		return err
	}
	return nil
}

func (r *UserRedisRepo) Sadds(ctx context.Context, key string, values []string) error {
	if err := r.redisClient.SAdd(ctx, key, values).Err(); err != nil {
		return err
	}
	return nil
}

func (r *UserRedisRepo) Srem(ctx context.Context, key string, value string) error {
	if err := r.redisClient.SRem(ctx, key, value).Err(); err != nil {
		return err
	}
	return nil
}

func (r *UserRedisRepo) SIsMember(ctx context.Context, key string, value string) (bool, error) {
	result := r.redisClient.SIsMember(ctx, key, value)
	if result.Err() != nil {
		return false, result.Err()
	}
	return result.Val(), nil
}
