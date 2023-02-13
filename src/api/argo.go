package api

import (
	"bytes"
	"fmt"
	"gepaplexx/git-workflows/logger"
	"gepaplexx/git-workflows/model"
	"gepaplexx/git-workflows/utils"
	"github.com/go-git/go-git/v5"
	"os/exec"
)

func UpdateArgoApplicationSet(c *model.Config, repo *git.Repository) {
	logger.Info("Updating ArgoCD Application")
	logger.Debug("Env: %s", c.Env())

	wt := checkout(repo, "main", false)

	if "main" == c.Env() {
		updateAllStages(c, wt)
	} else {
		filePath := fmt.Sprintf("%s/apps/env/%s/values.yaml", wt.Filesystem.Root(), c.Env())
		logger.Debug("Updating file: %s", filePath)
		updateImageTag(c, filePath)
	}
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
	logger.Debug(cmd.String())
	err := cmd.Run()
	utils.CheckIfError(err)
}
