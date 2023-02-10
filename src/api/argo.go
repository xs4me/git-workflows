package api

import (
	"bytes"
	"fmt"
	"gepaplexx/git-workflows/logger"
	"gepaplexx/git-workflows/model"
	"gepaplexx/git-workflows/utils"
	"github.com/go-git/go-git/v5"
	"os/exec"
	"strings"
)

func UpdateMultiDir(c *model.Config, repo *git.Repository) {
	logger.Info("Updating ArgoCD Application")
	env := strings.ReplaceAll(strings.ReplaceAll(c.Branch, "/", "-"), "_", "-")
	logger.Debug("Env: %s", env)
	c.Branch = "main"

	wt := checkout(c, repo)

	if "main" == env {
		updateAllStages(c, wt)
	} else {
		filePath := fmt.Sprintf("%s/apps/env/%s/values.yaml", wt.Filesystem.Root(), env)
		logger.Debug("Updating file: %s", filePath)
		updateImageTag(c, filePath)
	}
	commitAndPush(c, wt, repo, fmt.Sprintf("updated image tag to %s", c.ImageTag))
}

func UpdateMultiBranch(c *model.Config, repo *git.Repository) {
	wt := checkout(c, repo)
	filePath := fmt.Sprintf("%s/values.yaml", wt.Filesystem.Root())
	logger.Debug("Updating file: %s", filePath)
	updateImageTag(c, filePath)
	commitAndPush(c, wt, repo, fmt.Sprintf("updated image tag to %s", c.ImageTag))

}

func updateAllStages(c *model.Config, wt *git.Worktree) {

	for _, stage := range c.Stages {
		filePath := fmt.Sprintf("%s/apps/env/%s/values.yaml", wt.Filesystem.Root(), stage)
		logger.Debug("Updating file: %s", filePath)
		updateImageTag(c, filePath)
	}
}

func updateImageTag(c *model.Config, filePath string) {
	var out bytes.Buffer
	cmd := exec.Command("yq", "-i", fmt.Sprintf("with(%s; . = \"%s\" | . style=\"double\")", c.ImageTagLocation(), c.ImageTag), filePath)
	cmd.Stdout = &out
	cmd.Stderr = &out
	fmt.Println(cmd.String())
	err := cmd.Run()
	utils.CheckIfError(err)
}
