package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/hiennguyen9874/go-boilerplate-v2/config"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var (
	RootCmd = &cobra.Command{
		Use:   "go-base",
		Short: "A RESTful API boilerplate",
		Long:  `A RESTful API boilerplate with password less authentication.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		//	Run: func(cmd *cobra.Command, args []string) { },
	}
)

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-base.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	cfgViper, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	_, err = config.ParseConfig(cfgViper)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}
}
