package api

import (
	"bytes"
	"fmt"
	"gepaplexx/git-workflows/logger"
	"gepaplexx/git-workflows/model"
	"gepaplexx/git-workflows/utils"
	"github.com/go-git/go-git/v5"
	"github.com/otiai10/copy"
	"os/exec"
	"strings"
)

const (
	ELEMENT         = "{\"cluster\": \"%s\", \"url\": \"https://kubernetes.default.svc\", \"branch\": \"main\"}"
	SOURCE_SELECTOR = ".spec.generators[0].list.elements | map(select(.branch == \"%s\")) | .[0].cluster"
)

func UpdateArgoApplicationSet(c *model.Config, repo *git.Repository) {
	logger.Info("Updating ArgoCD Application")
	logger.Debug("Env: %s", c.Env)

	wt := checkout(repo, "main", false)

	if "main" == c.Env {
		updateAllStages(c, wt)
	} else {
		filePath := fmt.Sprintf("%s/apps/env/%s/values.yaml", wt.Filesystem.Root(), c.Env)
		logger.Debug("Updating file: %s", filePath)
		updateImageTag(c, filePath)
	}
	commitAndPush(c, wt, repo, fmt.Sprintf("updated image tag to %s", c.ImageTag))
}

func ArgoCreateEnvironment(c *model.Config, repo *git.Repository) {
	logger.Info("Creating entry in ArgoCD ApplicationSet")
	logger.Debug("Env: %s", c.Env)

	wt := checkout(repo, "main", false)
	filePath := fmt.Sprintf("%s/argocd/applicationset.yaml", wt.Filesystem.Root())
	addEnvironmentToApplicationSet(c, filePath)
	copyTemplateDir(wt, c, filePath)

	logger.Info("Copying ApplicationSet")
	err := copy.Copy(filePath, fmt.Sprintf("%s/appliationset.yaml", c.BaseDir))
	utils.CheckIfError(err)
	commitAndPush(c, wt, repo, fmt.Sprintf("Added ApplicationSet entry and templates for %s", c.Env))
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

func addEnvironmentToApplicationSet(c *model.Config, filePath string) {
	var out bytes.Buffer
	logger.Info("Adding %s to ApplicationSet %s", fmt.Sprintf(ELEMENT, c.Env), filePath)
	cmd := exec.Command("yq", "-i", fmt.Sprintf(".spec.generators[0].list.elements += %s", fmt.Sprintf(ELEMENT, c.Env)), filePath)
	cmd.Stdout = &out
	cmd.Stderr = &out
	logger.Debug(cmd.String())
	err := cmd.Run()
	utils.CheckIfError(err)
}

func copyTemplateDir(wt *git.Worktree, c *model.Config, applicationset string) {
	var out bytes.Buffer
	cmd := exec.Command("yq", fmt.Sprintf(SOURCE_SELECTOR, c.FromBranch), applicationset)
	cmd.Stdout = &out
	cmd.Stderr = &out
	logger.Debug(cmd.String())
	err := cmd.Run()
	utils.CheckIfError(err)

	sourceDir := fmt.Sprintf("%s/apps/env/%s", wt.Filesystem.Root(), strings.TrimRight(out.String(), "\n"))
	targetTemplateDir := fmt.Sprintf("%s/apps/env/%s", wt.Filesystem.Root(), c.Env)

	logger.Info("Copying %s to %s", sourceDir, targetTemplateDir)
	err = copy.Copy(sourceDir, targetTemplateDir)
	utils.CheckIfError(err)
}
