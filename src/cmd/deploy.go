package cmd

import (
	"gepaplexx/git-workflows/api"
	"gepaplexx/git-workflows/logger"
	"gepaplexx/git-workflows/model"
	"gepaplexx/git-workflows/utils"
	"github.com/spf13/cobra"
	"strings"
)

func init() {
	rootCmd.AddCommand(deployCmd)
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "handles the deployment accross multiple stages",
	Run: func(cmd *cobra.Command, args []string) {
		deploy(&Config)
	},
}

func deploy(c *model.Config) {
	sanitizeInput(c)
	if c.Development {
		developmentMode(c)
		c.FromBranch = "main"
		c.ToBranch = "prod"
	}
	deployPrerequisites(c)

	repo := api.CloneRepo(c, c.FromBranch, false)
	if !c.ResourcesOnly {
		wt, err := repo.Worktree()
		utils.CheckIfError(err)
		api.UpdateAllStages(c, wt, repo)
	}
	api.DeployFromTo(c, repo)
}

func deployPrerequisites(c *model.Config) {
	if c.GitUrl == "" {
		logger.Fatal("git url is missing")
	}
	if c.FromBranch == "" {
		logger.Fatal("from-branch is missing")
	}
	if c.ToBranch == "" {
		logger.Fatal("to-branch is missing")
	}
	if len(c.Stages) == 0 {
		logger.Fatal("stages are not configured")
	}
}

func sanitizeInput(c *model.Config) {
	c.FromBranch = strings.ReplaceAll(c.FromBranch, "[", "")
	c.FromBranch = strings.ReplaceAll(c.FromBranch, "]", "")
	c.ToBranch = strings.ReplaceAll(c.ToBranch, "[", "")
	c.ToBranch = strings.ReplaceAll(c.ToBranch, "]", "")
	c.Reponame = strings.ReplaceAll(c.Reponame, "[", "")
	c.Reponame = strings.ReplaceAll(c.Reponame, "]", "")
}
