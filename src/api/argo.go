package api

import (
	"fmt"
	"gepaplexx/git-workflows/logger"
	"gepaplexx/git-workflows/model"
	"gepaplexx/git-workflows/utils"
	"github.com/go-git/go-git/v5"
	"github.com/otiai10/copy"
	"os"
	"os/exec"
	"strings"
)

const (
	ELEMENT                 = "{\"cluster\": \"%s\", \"url\": \"https://kubernetes.default.svc\", \"branch\": \"main\"}"
	SOURCE_SELECTOR         = ".spec.generators[0].list.elements | map(select(.branch == \"%s\")) | .[0].cluster"
	APPLICATIONSET_LOCATION = "%s/argocd/applicationset.yaml"
	VALUES_LOCATION         = "%s/apps/env/%s/values.yaml"
	UPDATE_FORMAT           = "with(%s; . = \"%s\" | . style=\"double\")"
	ADD_FORMAT              = ".spec.generators[0].list.elements += %s"
	DELETE_FORMAT           = "del(.spec.generators[0].list.elements[] | select(.cluster == \"%s\"))"
	TEMPLATE_LOCATION       = "%s/apps/env/%s"
)

func UpdateArgoApplicationSet(c *model.Config, repo *git.Repository) {
	logger.Info("Updating ArgoCD Application")
	logger.Debug("Env: %s", c.Env)

	wt := checkout(repo, "main", false)

	if "main" == c.Env {
		updateAllStages(c, wt)
	} else {
		filePath := fmt.Sprintf(VALUES_LOCATION, wt.Filesystem.Root(), c.Env)
		logger.Debug("Updating file: %s", filePath)
		updateImageTag(c, filePath)
	}
	commitAndPush(c, wt, repo, fmt.Sprintf("updated image tag to %s", c.ImageTag))
}

func ArgoCreateEnvironment(c *model.Config, repo *git.Repository) {
	logger.Info("Creating entry in ArgoCD ApplicationSet")
	logger.Debug("Env: %s", c.Env)

	wt := checkout(repo, "main", false)
	filePath := fmt.Sprintf(APPLICATIONSET_LOCATION, wt.Filesystem.Root())
	addEnvironmentToApplicationSet(c, filePath)
	copyTemplateDir(wt, c, filePath)

	copyApplicationSet(c, filePath)
	commitAndPush(c, wt, repo, fmt.Sprintf("Added ApplicationSet entry and templates for %s", c.Env))
}

func DeleteArgoEnvironment(c *model.Config, repo *git.Repository) {
	logger.Info("Deleting entry in ArgoCD ApplicationSet")
	logger.Debug("Env: %s", c.Env)

	protectEnvironments(c)

	wt := checkout(repo, "main", false)
	filePath := fmt.Sprintf(APPLICATIONSET_LOCATION, wt.Filesystem.Root())
	removeEnvironmentFromApplicationSet(c, filePath)
	deleteTemplateDir(wt, c)
	copyApplicationSet(c, filePath)
	commitAndPush(c, wt, repo, fmt.Sprintf("Removed ApplicationSet entry and templates for %s", c.Env))
}

func protectEnvironments(c *model.Config) {
	for _, stage := range c.Stages {
		if stage == c.Branch && c.Force == false {
			logger.Fatal("%s is a predefined stage. It cannot be deleted via automation. Override this check by setting --force flag", stage)
		}
	}
}

func updateAllStages(c *model.Config, wt *git.Worktree) {

	for _, stage := range c.Stages {
		filePath := fmt.Sprintf(VALUES_LOCATION, wt.Filesystem.Root(), stage)
		logger.Debug("Updating file: %s", filePath)
		updateImageTag(c, filePath)
	}
}

func updateImageTag(c *model.Config, filepath string) {
	nodes := ParseYaml(filepath)
	tagNode, err := FindNode(nodes.Content[0], c.ImageTagLocation())
	utils.CheckIfError(err)
	tagNode.Value = c.ImageTag
	WriteYaml(nodes, filepath)
}

func addEnvironmentToApplicationSet(c *model.Config, path string) {
	logger.Info("Adding (%s, %s) to ApplicationSet %s", c.Env, "main", path)
	nodes := ParseYaml(path)
	envNode, err := FindNode(nodes.Content[0], "spec.generators.list.elements")
	utils.CheckIfError(err)
	envNode.Content = append(envNode.Content, NewEnvNode(c.Env, "main"))
	WriteYaml(nodes, path)
}

func copyTemplateDir(wt *git.Worktree, c *model.Config, applicationset string) {
	fromBranch := strings.ReplaceAll(strings.ReplaceAll(c.FromBranch, "/", "-"), "_", "-")
	cmd := exec.Command("yq", fmt.Sprintf(SOURCE_SELECTOR, fromBranch), applicationset)
	res := execute(cmd)

	sourceDir := fmt.Sprintf(TEMPLATE_LOCATION, wt.Filesystem.Root(), res)
	targetTemplateDir := fmt.Sprintf(TEMPLATE_LOCATION, wt.Filesystem.Root(), c.Env)

	logger.Info("Copying %s to %s", sourceDir, targetTemplateDir)
	err := copy.Copy(sourceDir, targetTemplateDir)
	utils.CheckIfError(err)
}

func deleteTemplateDir(wt *git.Worktree, c *model.Config) {
	targetTemplateDir := fmt.Sprintf(TEMPLATE_LOCATION, wt.Filesystem.Root(), c.Env)

	logger.Info("Deleting %s", targetTemplateDir)
	err := os.RemoveAll(targetTemplateDir)
	utils.CheckIfError(err)
}

func removeEnvironmentFromApplicationSet(c *model.Config, path string) {
	logger.Info("Removing %s from ApplicationSet %s", c.Env, path)
	cmd := exec.Command("yq", "-i", fmt.Sprintf(DELETE_FORMAT, c.Env), path)
	_ = execute(cmd)
}

func copyApplicationSet(c *model.Config, filePath string) {
	logger.Info("Copying ApplicationSet")
	err := copy.Copy(filePath, fmt.Sprintf("%s/application.yaml", c.BaseDir))
	utils.CheckIfError(err)
}
