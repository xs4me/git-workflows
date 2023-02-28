package api

import (
	"bytes"
	"errors"
	"fmt"
	"gepaplexx/git-workflows/utils"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"strings"
)

const (
	AppsetEnvPath = ""
)

func ParseYaml(filepath string) yaml.Node {
	file, err := os.Open(filepath)
	utils.CheckIfError(err)

	defer func(file *os.File) {
		err := file.Close()
		utils.CheckIfError(err)
	}(file)

	by, err := io.ReadAll(file)
	utils.CheckIfError(err)

	var node yaml.Node
	err = yaml.Unmarshal(by, &node)
	utils.CheckIfError(err)
	return node
}
func WriteYaml(nodes yaml.Node, filepath string) {
	var b bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(2)
	err := yamlEncoder.Encode(&nodes)
	utils.CheckIfError(err)

	err = os.WriteFile(filepath, b.Bytes(), 0664)
	utils.CheckIfError(err)
}

func DeleteEnvFromApplicationset(rootNode *yaml.Node, env string) error {
	envNode, err := FindNode(rootNode, AppsetEnvPath)
	utils.CheckIfError(err)
	idxToDelete := -1
	for idx, n := range envNode.Content {
		for i := 0; i < len(n.Content)-1 && idxToDelete == -1; i += 2 {
			if n.Content[i].Value == "cluster" && n.Content[i+1].Value == env {
				idxToDelete = idx
			}
		}
	}

	if idxToDelete == -1 {
		return errors.New(fmt.Sprintf("could not find environment '%s'", env))
	}

	removeIndex(envNode.Content, idxToDelete)
	return nil
}

func FindClusterWithBranch(rootNode *yaml.Node, branch string) (string, error) {
	envNode, err := FindNode(rootNode, AppsetEnvPath)
	if err != nil {
		panic(err)
	}

	for _, n := range envNode.Content {
		for i := 0; i < len(n.Content)-1; i += 2 {
			if n.Content[i].Value == "branch" && n.Content[i+1].Value == branch {
				return extractCluster(n)
			}
		}
	}

	return "", errors.New(fmt.Sprintf("no branch with name '%s' found", branch))
}

func NewEnvNode(env, branch string) *yaml.Node {
	newEnvNode := yaml.Node{
		Kind: yaml.MappingNode,
	}
	newEnvNode.Content = append(newEnvNode.Content, newScalarNode("cluster"))
	newEnvNode.Content = append(newEnvNode.Content, newScalarNode(env))
	newEnvNode.Content = append(newEnvNode.Content, newScalarNode("branch"))
	newEnvNode.Content = append(newEnvNode.Content, newScalarNode(branch))
	newEnvNode.Content = append(newEnvNode.Content, newScalarNode("url"))
	newEnvNode.Content = append(newEnvNode.Content, newScalarNode("https://kubernetes.default.svc"))
	return &newEnvNode
}

func FindNode(rootNode *yaml.Node, lookingFor string) (*yaml.Node, error) {
	current := ""
	return find(rootNode, lookingFor, &current)
}

func find(node *yaml.Node, lookingFor string, current *string) (*yaml.Node, error) {
	switch node.Kind {
	case yaml.MappingNode:
		{
			if found := handleMappingNode(node, lookingFor, current); found != nil {
				return found, nil
			}
		}
	case yaml.SequenceNode:
		{
			if found := handleSequenceNode(node, lookingFor, current); found != nil {
				return found, nil
			}
		}
	case yaml.ScalarNode:
		{
			if node.Value == lookingFor {
				return node, nil
			}
		}
	}
	return nil, errors.New("element not found")
}

func handleMappingNode(node *yaml.Node, lookingFor string, current *string) *yaml.Node {
	for i := 0; i < len(node.Content)-1; i += 2 {
		appendIfValid(current, node.Content[i].Value, lookingFor)
		if *current == lookingFor {
			return node.Content[i+1]
		}
		found, err := find(node.Content[i+1], lookingFor, current)
		if err == nil {
			return found
		}
	}

	return nil
}

func handleSequenceNode(node *yaml.Node, lookingFor string, current *string) *yaml.Node {
	for _, n := range node.Content {
		found, err := find(n, lookingFor, current)
		if err == nil {
			return found
		}
	}
	return nil
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

func newScalarNode(value string) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: value,
	}
}

func removeIndex(s []*yaml.Node, index int) []*yaml.Node {
	return append(s[:index], s[index+1:]...)
}

func extractCluster(node *yaml.Node) (string, error) {
	for i := 0; i < len(node.Content)-1; i += 2 {
		if node.Content[i].Value == "cluster" {
			return node.Content[i+1].Value, nil
		}
	}

	return "", errors.New("no cluster tag found")
}
