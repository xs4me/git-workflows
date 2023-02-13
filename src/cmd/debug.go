package cmd

import (
	"gepaplexx/git-workflows/logger"
	"gepaplexx/git-workflows/model"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(debugCmd)
}

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "prints current configuration of the application",
	Run: func(cmd *cobra.Command, args []string) {
		printConfig(&Config)
	},
}

func printConfig(c *model.Config) {
	if c.Development {
		logger.Debug("Development mode enabled. Using local configuration.")
		developmentMode(c)
	}

	logger.Info("Config: %+v", Config)
	logger.Info("Application Version: %s", Version)
}
