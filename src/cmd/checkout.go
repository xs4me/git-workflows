package cmd

import (
	"gepaplexx/git-workflows/api"
	"gepaplexx/git-workflows/logger"
	"gepaplexx/git-workflows/model"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(checkoutCmd)
}

var checkoutCmd = &cobra.Command{
	Use:   "checkout",
	Short: "handles the checkout stage of the workflow",
	Run: func(cmd *cobra.Command, args []string) {
		checkout(&Config)
	},
}

func checkout(c *model.Config) {
	if c.Development {
		logger.Debug("Development mode enabled. Using local configuration.")
		developmentMode(c)
	}

	logger.Info("checkout %s\n", c.GitUrl)
	logger.Debug("Checking prerequisites...")
	prerequisites(c)

	api.CloneAppRepo(c)
	api.ExtractGitInformation(c)
}

func prerequisites(c *model.Config) {
	if c.GitUrl == "" {
		logger.Error("GitUrl is empty. Please provide a valid git url.")
		os.Exit(1)
	}
	if c.Reponame == "" {
		logger.Error("Repository name is empty. Please provide a valid repository name.")
		os.Exit(1)
	}
}
