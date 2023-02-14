package cmd

import (
	"gepaplexx/git-workflows/api"
	"gepaplexx/git-workflows/model"
	"github.com/spf13/cobra"
	"strings"
)

var createCmd = &cobra.Command{
	Use:   "argo-create",
	Short: "Creates a new argocd application on branch creation and updates git repository accordingly.",
	Run: func(cmd *cobra.Command, args []string) {
		createArgo(&Config)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func createArgo(c *model.Config) {
	c.Env = strings.ReplaceAll(strings.ReplaceAll(c.Branch, "/", "-"), "_", "-")
	if c.Development {
		developmentMode(c)
		c.Branch = "feature/new-branch"
		c.Env = strings.ReplaceAll(strings.ReplaceAll(c.Branch, "/", "-"), "_", "-")
		c.FromBranch = "qa"
	}
	argoCreatePrerequisites(c)
	repo := api.CloneRepo(c, "main", false)
	api.ArgoCreateEnvironment(c, repo)
}

func argoCreatePrerequisites(c *model.Config) {
	prerequisites(c)
}
