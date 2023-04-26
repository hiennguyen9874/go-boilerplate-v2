package cmd

import (
	"context"

	"github.com/hiennguyen9874/go-boilerplate-v2/config"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/distributor"
	userDistributor "github.com/hiennguyen9874/go-boilerplate-v2/internal/users/distributor"
	userRepository "github.com/hiennguyen9874/go-boilerplate-v2/internal/users/repository"
	userUseCase "github.com/hiennguyen9874/go-boilerplate-v2/internal/users/usecase"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/db/postgres"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/db/redis"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/logger"
	"github.com/spf13/cobra"
)

var initDataCmd = &cobra.Command{
	Use:   "initdata",
	Short: "Init data",
	Long:  "Init data",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetCfg()

		appLogger := logger.NewApiLogger(cfg)
		appLogger.InitLogger()
		appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

		psqlDB, err := postgres.NewPsqlDB(cfg)
		if err != nil {
			appLogger.Fatalf("Postgresql init: %s", err)
		} else {
			appLogger.Infof("Postgres connected")
		}

		redisClient := redis.NewRedis(cfg)

		taskRedisClient := distributor.NewRedisClient(cfg)

		// Repository
		userPgRepo := userRepository.CreateUserPgRepository(psqlDB)
		userRedisRepo := userRepository.CreateUserRedisRepository(redisClient)

		// Distributor
		userRedisTaskDistributor := userDistributor.NewUserRedisTaskDistributor(taskRedisClient, cfg, appLogger)

		// UseCase
		userUC := userUseCase.CreateUserUseCaseI(userPgRepo, userRedisRepo, userRedisTaskDistributor, cfg, appLogger)

		// Create super user if not exists
		isCreated, _ := userUC.CreateSuperUserIfNotExist(context.Background())

		if !isCreated {
			appLogger.Info("Super user is exists, skip create")
		} else {
			appLogger.Info("Created super user")
		}
	},
}

func init() {
	RootCmd.AddCommand(initDataCmd)
}
