package cmd

import (
	"context"

	"github.com/hiennguyen9874/go-boilerplate-v2/config"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/distributor"
	"github.com/hiennguyen9874/go-boilerplate-v2/internal/server"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/db/postgres"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/db/redis"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/logger"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start http server with configured api",
	Long:  `Starts a http server and serves the configured api`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		cfg := config.GetCfg()

		appLogger := logger.NewApiLogger(cfg)
		appLogger.InitLogger()
		appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

		psqlClient, err := postgres.NewPsqlClient(cfg)
		if err != nil {
			appLogger.Fatalf("Postgresql init: %s", err)
		} else {
			appLogger.Infof("Postgres connected")
		}

		if cfg.Server.MigrateOnStart {
			err = Migrate(ctx, psqlClient)

			if err != nil {
				appLogger.Info("Can not migrate data", err)
			} else {
				appLogger.Info("Data migrated")
			}
		}

		redisClient := redis.NewRedis(cfg)

		taskRedisClient := distributor.NewRedisClient(cfg)

		server, err := server.NewServer(cfg, psqlClient, redisClient, taskRedisClient, appLogger)
		if err != nil {
			appLogger.Fatal(err)
		}
		server.Start()
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}
