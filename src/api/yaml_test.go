package api

import (
	"gopkg.in/yaml.v3"
	"testing"
)

// region TESTDATA
var inputYamlV1 = `demo-microservice:
  replicaCount: 1
  image:
    name: ghcr.io/gepaplexx/demo-microservice
    tag: "413b395"
  ports:
    - name: http
      containerPort: 8080
      protocol: TCP
`

var inputYamlV2 = `replicaCount: 1
image:
  name: ghcr.io/gepaplexx/demo-microservice
  tag: "413b395"
ports:
  - name: http
    containerPort: 8080
    protocol: TCP
`

// endregion

type findNodeTest struct {
	inputYaml, expectedValue, yamlPath string
}

var findNodeTests = []findNodeTest{
	{inputYamlV1, "413b395", "demo-microservice.image.tag"},
	{inputYamlV2, "413b395", "image.tag"},
	{inputYamlV2, "8080", "ports.containerPort"},
}

func TestFindNode(t *testing.T) {
	for _, test := range findNodeTests {
		nodes, err := unmarshal(test.inputYaml)
		if err != nil {
			t.Fatalf("failed to unmarshal: '%s'", test.inputYaml)
		}

		node, err := FindNode(nodes.Content[0], test.yamlPath)
		if err != nil {
			t.Fatalf("node %s not found", test.yamlPath)
		}

		if node.Value != test.expectedValue {
			t.Fatalf("expected value '%s' did not match '%s'", test.expectedValue, node.Value)
		}
	}
}

func TestFindNodeNotFound(t *testing.T) {
	yamlPath := "did.not.exist"
	nodes, err := unmarshal(inputYamlV1)
	if err != nil {
		t.Fatalf("failed to unmarshal: '%s'", inputYamlV1)
	}

	node, err := FindNode(nodes.Content[0], yamlPath)
	if err == nil {
		t.Fatal("error shouldn't be nil")
	}

	if err.Error() != "element not found" {
		t.Fatal("unexpected error message")
	}

	if node != nil {
		t.Fatal("node shouldn't be found, should be nil")
	}
}

// TODO TESTS für DeleteEnvFromApplicationset, FindClusterWithBranch

func unmarshal(data string) (yaml.Node, error) {
	var nodes yaml.Node
	err := yaml.Unmarshal([]byte(data), &nodes)
	if err != nil {
		return yaml.Node{}, err
	}

	return nodes, nil
}
