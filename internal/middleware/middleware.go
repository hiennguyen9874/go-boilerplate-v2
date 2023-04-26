package middleware

import (
	"github.com/hiennguyen9874/go-boilerplate/config"
	"github.com/hiennguyen9874/go-boilerplate/internal/users"
	"github.com/hiennguyen9874/go-boilerplate/pkg/logger"
)

type MiddlewareManager struct {
	cfg     *config.Config
	logger  logger.Logger
	usersUC users.UserUseCaseI
}

func CreateMiddlewareManager(cfg *config.Config, logger logger.Logger, usersUC users.UserUseCaseI) *MiddlewareManager {
	return &MiddlewareManager{
		cfg:     cfg,
		logger:  logger,
		usersUC: usersUC,
	}
}
