package users

import (
	"github.com/hiennguyen9874/go-boilerplate/internal"
	"github.com/hiennguyen9874/go-boilerplate/internal/models"
)

type UserRedisRepository interface {
	internal.RedisRepository[models.User]
}
