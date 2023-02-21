package api

import (
	"bytes"
	"fmt"
	"gepaplexx/git-workflows/logger"
	"gepaplexx/git-workflows/model"
	"gepaplexx/git-workflows/utils"
	"github.com/go-git/go-git/v5"
	"github.com/otiai10/copy"
	"gopkg.in/yaml.v3"
	"io"
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
	nodes := parseYaml(filepath)
	updated := updateVal(nodes.Content[0], c.ImageTagLocation(), c.ImageTag)
	if !updated {
		logger.Fatal("no update happend: %s not found", c.ImageTagLocation())
	}

	writeUpdatedYaml(nodes, filepath)
}

func parseYaml(filepath string) yaml.Node {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	by, err := io.ReadAll(file)
	utils.CheckIfError(err)

	var node yaml.Node
	err = yaml.Unmarshal(by, &node)
	utils.CheckIfError(err)
	return node
}

func writeUpdatedYaml(nodes yaml.Node, filepath string) {
	var b bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(2)
	err := yamlEncoder.Encode(&nodes)
	utils.CheckIfError(err)

	err = os.WriteFile(filepath, b.Bytes(), 0664)
	utils.CheckIfError(err)
}
func updateVal(node *yaml.Node, updatePath string, newVal string) bool {
	current := ""
	found := false
	update(node, &current, updatePath, newVal, &found)
	return found
}

func update(node *yaml.Node, current *string, lookingFor string, newVal string, found *bool) {
	if *found {
		return
	}
	if node.Kind == yaml.SequenceNode {
		for _, child := range node.Content {
			update(child, current, lookingFor, newVal, found)
		}
	} else if node.Kind == yaml.MappingNode {
		for i := 0; i < len(node.Content); i += 2 {
			key := node.Content[i]
			value := node.Content[i+1]
			appendIfValid(current, key.Value, lookingFor)
			if *current == lookingFor {
				node.Content[i+1].Value = newVal
				*found = true
				return
			}
			update(key, current, lookingFor, newVal, found)
			update(value, current, lookingFor, newVal, found)

		}
	} else {
		if *current == lookingFor {
			node.Value = newVal
			*found = true
		}

		appendIfValid(current, node.Value, lookingFor)
	}
}

func appendIfValid(current *string, appendix string, lookingFor string) {
	newCurrent := *current
	if *current != "" {
		newCurrent += "."
	}

	newCurrent += appendix
	if !strings.HasPrefix(lookingFor, newCurrent) {
		if *current != "" {
			newCurrent = strings.TrimSuffix(newCurrent, "."+appendix)
		} else {
			newCurrent = ""
		}
	}

	*current = newCurrent
}

func addEnvironmentToApplicationSet(c *model.Config, path string) {
	logger.Info("Adding %s to ApplicationSet %s", fmt.Sprintf(ELEMENT, c.Env), path)
	cmd := exec.Command("yq", "-i", fmt.Sprintf(ADD_FORMAT, fmt.Sprintf(ELEMENT, c.Env)), path)
	_ = execute(cmd)
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
	err := copy.Copy(filePath, fmt.Sprintf("%s/appliationset.yaml", c.BaseDir))
	utils.CheckIfError(err)
}
