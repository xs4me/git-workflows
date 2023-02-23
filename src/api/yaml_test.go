package api

import (
	"bytes"
	"gopkg.in/yaml.v3"
	"strings"
	"testing"
)

// region TESTDATA
var inputYamlTagV1 = `demo-microservice:
  replicaCount: 1
  image:
    name: ghcr.io/gepaplexx/demo-microservice
    tag: "413b395"
  ports:
    - name: http
      containerPort: 8080
      protocol: TCP
`
var expectedYamlTagV1 = `demo-microservice:
  replicaCount: 1
  image:
    name: ghcr.io/gepaplexx/demo-microservice
    tag: "1234567"
  ports:
    - name: http
      containerPort: 8080
      protocol: TCP
`
var inputYamlTagV2 = `replicaCount: 1
image:
  name: ghcr.io/gepaplexx/demo-microservice
  tag: "413b395"
ports:
  - name: http
    containerPort: 8080
    protocol: TCP
`
var expectedYamlTagV2 = `replicaCount: 1
image:
  name: ghcr.io/gepaplexx/demo-microservice
  tag: "1234567"
ports:
  - name: http
    containerPort: 8080
    protocol: TCP
`

var inputYamlPort = `replicaCount: 1
image:
  name: ghcr.io/gepaplexx/demo-microservice
  tag: "413b395"
ports:
  - name: http
    containerPort: 8080
    protocol: TCP
`
var expectedYamlPort = `replicaCount: 1
image:
  name: ghcr.io/gepaplexx/demo-microservice
  tag: "413b395"
ports:
  - name: http
    containerPort: 1234567
    protocol: TCP
`

// endregion

type updateValTest struct {
	inputYaml, expectedYaml, yamlPath string
}

var updateValTests = []updateValTest{
	{inputYamlTagV1, expectedYamlTagV1, "demo-microservice.image.tag"},
	{inputYamlTagV2, expectedYamlTagV2, "image.tag"},
	{inputYamlPort, expectedYamlPort, "ports.containerPort"},
}

func TestUpdateVal(t *testing.T) {
	newVal := "1234567"
	for _, test := range updateValTests {
		nodes, err := unmarshal(test.inputYaml)
		if err != nil {
			t.Fatalf("failed to unmarshal: '%s'", test.inputYaml)
		}

		updated := updateVal(nodes.Content[0], test.yamlPath, newVal)
		if !updated {
			t.Fatal("yaml was not updated")
		}

		updatedYaml, err := marshal(nodes)
		if err != nil {
			t.Fatal("failed to marshal data")
		}

		if strings.Compare(test.expectedYaml, updatedYaml) != 0 {
			t.Fatalf("updated yaml was not equal to expected yaml: expected: '%s', current: '%s'", test.expectedYaml, updatedYaml)
		}
	}
}

func TestUpdateValNotFound(t *testing.T) {
	newVal := "doesNotMatter"
	yamlPath := "did.not.exist"
	nodes, err := unmarshal(inputYamlTagV1)
	if err != nil {
		t.Fatalf("failed to unmarshal: '%s'", inputYamlTagV1)
	}

	updated := updateVal(nodes.Content[0], yamlPath, newVal)
	if updated {
		t.Fatalf("unexpected update, path: '%s', input: '%s'", inputYamlTagV1, yamlPath)
	}

	updatedYaml, err := marshal(nodes)
	if err != nil {
		t.Fatal("failed to marshal data")
	}

	if strings.Compare(inputYamlTagV1, updatedYaml) != 0 {
		t.Fatalf("input did not match output, shouldn't be updated. input: '%s', output: '%s'", inputYamlTagV1, updatedYaml)
	}
}

func unmarshal(data string) (yaml.Node, error) {
	var nodes yaml.Node
	err := yaml.Unmarshal([]byte(data), &nodes)
	if err != nil {
		return yaml.Node{}, err
	}

	return nodes, nil
}

func marshal(nodes yaml.Node) (string, error) {
	var b bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(2)
	err := yamlEncoder.Encode(&nodes)
	if err != nil {
		return "", err
	}

	return string(b.Bytes()), nil
}
