package cmd

import (
	"gepaplexx/git-workflows/api"
	"gepaplexx/git-workflows/logger"
	"gepaplexx/git-workflows/model"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "argo-update",
	Short: "Updates argocd application in infrastructure repository to handle deployments",
	Run: func(cmd *cobra.Command, args []string) {
		updateArgo(&Config)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func updateArgo(c *model.Config) {
	checkArgoprequesites(c)
	if c.Development {
		logger.Debug("Development mode enabled. Using local configuration.")
		developmentMode(c)
	}
	repo := api.CloneRepo(c, false)
	if c.LegacyBehavior {
		api.UpdateMultiBranch(c, repo)
	}
	if !c.LegacyBehavior {
		api.UpdateMultiDir(c, repo)
	}
}

func checkArgoprequesites(c *model.Config) {
	if c.ImageTag == "" {
		logger.Fatal("Image tag must be set")
	}

	if c.Branch == "" {
		logger.Fatal("Branch must be set")
	}

}
