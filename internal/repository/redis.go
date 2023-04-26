package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/hiennguyen9874/go-boilerplate/internal"
	"github.com/hiennguyen9874/go-boilerplate/pkg/httpErrors"
	"github.com/redis/go-redis/v9"
)

type RedisRepo[M any] struct {
	RedisClient *redis.Client
}

func CreateRedisRepo[M any](redisClient *redis.Client) RedisRepo[M] {
	return RedisRepo[M]{RedisClient: redisClient}
}

func CreateRedisRepository[M any](redisClient *redis.Client) internal.RedisRepository[M] {
	return &RedisRepo[M]{RedisClient: redisClient}
}

func (r *RedisRepo[M]) Create(ctx context.Context, key string, exp *M, seconds int) error {
	objBytes, err := json.Marshal(exp)
	if err != nil {
		return httpErrors.ErrJson(err)
	}

	if err = r.RedisClient.Set(ctx, key, objBytes, time.Second*time.Duration(seconds)).Err(); err != nil {
		// TODO: Using httpErrors
		return err
	}
	return nil
}

func (r *RedisRepo[M]) Get(ctx context.Context, key string) (*M, error) {
	objBytes, err := r.RedisClient.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}

	var obj M

	if err = json.Unmarshal(objBytes, &obj); err != nil {
		return nil, httpErrors.ErrJson(err)
	}

	return &obj, nil
}

func (r *RedisRepo[M]) Delete(ctx context.Context, key string) error {
	if err := r.RedisClient.Del(ctx, key).Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		// TODO: Using httpErrors
		return err
	}
	return nil
}

func (r *RedisRepo[M]) Sadd(ctx context.Context, key string, value string) error {
	if err := r.RedisClient.SAdd(ctx, key, value).Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisRepo[M]) Sadds(ctx context.Context, key string, values []string) error {
	if err := r.RedisClient.SAdd(ctx, key, values).Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisRepo[M]) Srem(ctx context.Context, key string, value string) error {
	if err := r.RedisClient.SRem(ctx, key, value).Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisRepo[M]) SIsMember(ctx context.Context, key string, value string) (bool, error) {
	result := r.RedisClient.SIsMember(ctx, key, value)
	if result.Err() != nil {
		return false, result.Err()
	}
	return result.Val(), nil
}

// func (r *RedisRepo[M]) SMembers(ctx context.Context, key string) ([]string, error) {
// 	result := r.RedisClient.SPop(ctx, key)
// 	if result.Err() != nil {
// 		return nil, result.Err()
// 	}
// 	return result.Val(), nil
// }
