package api

import (
	"gepaplexx/git-workflows/logger"
	"gepaplexx/git-workflows/model"
	"github.com/go-git/go-git/v5"
	"strings"
)

func UpdateMultiDir(c *model.Config, repo *git.Repository) {
	logger.Info("Updating ArgoCD Application")
	env := strings.ReplaceAll(strings.ReplaceAll(c.Branch, "/", "-"), "_", "-")
	logger.Debug("Env: %s", env)
	c.Branch = "main"

	_ = checkout(c, repo)

}

func UpdateMultiBranch(c *model.Config, repo *git.Repository) {
	logger.Fatal("Not implemented yet")
}
