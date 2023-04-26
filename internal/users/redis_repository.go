package users

import (
	"github.com/hiennguyen9874/go-boilerplate-v2/internal"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/models"
)

type UserRedisRepository interface {
	internal.RedisRepository[models.User]
}
