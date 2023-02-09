package api

import (
	"fmt"
	"gepaplexx/git-workflows/logger"
	"gepaplexx/git-workflows/model"
	"gepaplexx/git-workflows/utils"
	"github.com/go-git/go-git/v5"
	"gopkg.in/yaml.v3"
	"os"
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
	}
	commitAndPush(c, wt, repo, fmt.Sprintf("Update image tag to %s", c.ImageTag))
}

func UpdateMultiBranch(c *model.Config, repo *git.Repository) {
	logger.Fatal("Not implemented yet")
}

func updateAllStages(c *model.Config, wt *git.Worktree) {

	for _, stage := range c.Stages {
		filePath := fmt.Sprintf("%s/apps/env/%s/values.yaml", wt.Filesystem.Root(), stage)
		logger.Debug("Updating file: %s", filePath)
		updateImageTag(c, filePath)
	}
}

func updateImageTag(c *model.Config, filePath string) {
	values, err := os.ReadFile(filePath)
	utils.CheckIfError(err)

	data := make(map[string]interface{})
	err = yaml.Unmarshal(values, &data)

	findAndUpdateValue(data, c.ImageTag, strings.Split(c.ImageTagLocation(), ".")...)
	utils.CheckIfError(err)

	out, err := yaml.Marshal(data)
	utils.CheckIfError(err)
	err = os.WriteFile(filePath, out, 0644)
	utils.CheckIfError(err)
}

func findAndUpdateValue(m map[string]any, value string, keys ...string) (rval any) {
	var ok bool

	if len(keys) == 0 { // degenerate input
		logger.Fatal("NestedMapLookup needs at least one key")
	}
	if rval, ok = m[keys[0]]; !ok {
		logger.Fatal("key not found: %s", keys[0])
	}
	if len(keys) == 1 { // we've reached the final key
		m[keys[0]] = value
		return rval
	}
	if m, ok = rval.(map[string]any); !ok {
		logger.Fatal("malformed structure at %#v", rval)
	} else { // 1+ more keys
		return findAndUpdateValue(m, value, keys[1:]...)
	}
	return nil
}
