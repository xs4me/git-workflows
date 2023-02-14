package cmd

import (
	"gepaplexx/git-workflows/api"
	"gepaplexx/git-workflows/logger"
	"gepaplexx/git-workflows/model"
	"github.com/spf13/cobra"
	"strings"
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
	c.Env = strings.ReplaceAll(strings.ReplaceAll(c.Branch, "/", "-"), "_", "-")
	if c.Development {
		developmentMode(c)
		c.Branch = "main"
		c.Env = strings.ReplaceAll(strings.ReplaceAll(c.Branch, "/", "-"), "_", "-")
	}
	argoUpdatePrerequisites(c)
	repo := api.CloneRepo(c, "main", false)
	api.UpdateArgoApplicationSet(c, repo)
}

func argoUpdatePrerequisites(c *model.Config) {
	prerequisites(c)

	if c.ImageTag == "" {
		logger.Fatal("Image tag must be set")
	}
}
