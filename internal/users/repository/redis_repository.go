package repository

import (
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/models"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/repository"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/users"
	"github.com/redis/go-redis/v9"
)

type UserRedisRepo struct {
	repository.RedisRepo[models.User]
}

func CreateUserRedisRepository(redisClient *redis.Client) users.UserRedisRepository {
	return &UserRedisRepo{
		RedisRepo: repository.RedisRepo[models.User](repository.CreateRedisRepo[models.User](redisClient)),
	}
}
