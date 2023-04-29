package cmd

import (
	"context"

	"github.com/hiennguyen9874/go-boilerplate-v2/config"
	"github.com/hiennguyen9874/go-boilerplate-v2/ent"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/db/postgres"
	"github.com/hiennguyen9874/go-boilerplate-v2/pkg/logger"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate data",
	Long:  "Migrate data",
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

		err = Migrate(ctx, psqlClient)

		if err != nil {
			appLogger.Info("Can not migrate data", err)
		} else {
			appLogger.Info("Data migrated")
		}
	},
}

func Migrate(ctx context.Context, client *ent.Client) error {
	err := client.Schema.Create(ctx)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	RootCmd.AddCommand(migrateCmd)
}
