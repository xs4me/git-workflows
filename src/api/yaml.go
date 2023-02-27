package api

import (
	"bytes"
	"gepaplexx/git-workflows/utils"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"strings"
)

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

func parseYaml(filepath string) yaml.Node {
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
