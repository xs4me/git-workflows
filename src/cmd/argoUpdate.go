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
	logger.Debug("%+v", repo)
}

func checkArgoprequesites(c *model.Config) {

}
