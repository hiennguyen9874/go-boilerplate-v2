package cmd

import (
	"github.com/hiennguyen9874/go-boilerplate/config"
	"github.com/hiennguyen9874/go-boilerplate/internal/worker"
	"github.com/hiennguyen9874/go-boilerplate/pkg/logger"
	"github.com/spf13/cobra"
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "worker",
	Long:  "worker",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetCfg()

		appLogger := logger.NewApiLogger(cfg)
		appLogger.InitLogger()
		appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

		server, err := worker.NewTaskProcessor(cfg, appLogger)
		if err != nil {
			appLogger.Fatal(err)
		}
		server.Start() //nolint:errcheck
	},
}

func init() {
	RootCmd.AddCommand(workerCmd)
}
