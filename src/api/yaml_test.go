package api

import (
	"bytes"
	"fmt"
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

var inputYamlBaukastenV1 = `apiVersion: core.oam.dev/v1beta1
kind: Application
metadata:
  name: demo-microservice
spec:
  components:
    - name: demo-microservice
      type: deployment
      properties:
        image: "ghcr.io/gepaplexx/demo-microservice"
        tag: "d821097"
        ports:
          - port: 8080
            expose: true
`

var inputYamlBaukastenV2 = `apiVersion: core.oam.dev/v1beta1
kind: Application
metadata:
  name: demo-microservice
spec:
  components:
    - name: mega-backend
      type: deployment
      properties:
        image: "ghcr.io/gepaplexx/mega-backend"
        tag: "d821097"
    - name: mega-frontend
      type: deployment
      properties:
        image: "ghcr.io/gepaplexx/mega-frontend"
        tag: "1234567"
`

var applicationSetYaml = `apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: demo-microservice
  namespace: gepaplexx-cicd-tools
spec:
  generators:
    - list:
        elements:
          - cluster: dev
            url: https://kubernetes.default.svc
            branch: dev
          - cluster: qa
            url: https://kubernetes.default.svc
            branch: qa
          - cluster: main
            url: https://kubernetes.default.svc
            branch: main
          - cluster: feature-xy
            url: https://kubernetes.default.svc
            branch: main
  template:
    metadata:
      name: "demo-microservice-{{cluster}}"
    spec:
      project: demo-microservice
      source:
        repoURL: git@github.com:gepaplexx-demos/demo-microservice-ci.git
        targetRevision: "{{ branch }}"
        path: apps/env/{{ cluster }}
      destination:
        server: "{{url}}"
        namespace: "demo-microservice-{{cluster}}"
`

var expectedAppSetFirst = `apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: demo-microservice
  namespace: gepaplexx-cicd-tools
spec:
  generators:
    - list:
        elements:
          - cluster: qa
            url: https://kubernetes.default.svc
            branch: qa
          - cluster: main
            url: https://kubernetes.default.svc
            branch: main
          - cluster: feature-xy
            url: https://kubernetes.default.svc
            branch: main
  template:
    metadata:
      name: "demo-microservice-{{cluster}}"
    spec:
      project: demo-microservice
      source:
        repoURL: git@github.com:gepaplexx-demos/demo-microservice-ci.git
        targetRevision: "{{ branch }}"
        path: apps/env/{{ cluster }}
      destination:
        server: "{{url}}"
        namespace: "demo-microservice-{{cluster}}"
`

var expectedAppSetMiddle = `apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: demo-microservice
  namespace: gepaplexx-cicd-tools
spec:
  generators:
    - list:
        elements:
          - cluster: dev
            url: https://kubernetes.default.svc
            branch: dev
          - cluster: main
            url: https://kubernetes.default.svc
            branch: main
          - cluster: feature-xy
            url: https://kubernetes.default.svc
            branch: main
  template:
    metadata:
      name: "demo-microservice-{{cluster}}"
    spec:
      project: demo-microservice
      source:
        repoURL: git@github.com:gepaplexx-demos/demo-microservice-ci.git
        targetRevision: "{{ branch }}"
        path: apps/env/{{ cluster }}
      destination:
        server: "{{url}}"
        namespace: "demo-microservice-{{cluster}}"
`

var expectedAppSetLast = `apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: demo-microservice
  namespace: gepaplexx-cicd-tools
spec:
  generators:
    - list:
        elements:
          - cluster: dev
            url: https://kubernetes.default.svc
            branch: dev
          - cluster: qa
            url: https://kubernetes.default.svc
            branch: qa
          - cluster: main
            url: https://kubernetes.default.svc
            branch: main
  template:
    metadata:
      name: "demo-microservice-{{cluster}}"
    spec:
      project: demo-microservice
      source:
        repoURL: git@github.com:gepaplexx-demos/demo-microservice-ci.git
        targetRevision: "{{ branch }}"
        path: apps/env/{{ cluster }}
      destination:
        server: "{{url}}"
        namespace: "demo-microservice-{{cluster}}"
`

// endregion

// region FindNode Tests
type findNodeTest struct {
	inputYaml, expectedValue, yamlPath string
}

var findNodeTests = []findNodeTest{
	{inputYamlV1, "413b395", "demo-microservice.image.tag"},
	{inputYamlV2, "413b395", "image.tag"},
	{inputYamlV2, "8080", "ports.containerPort"},
	{inputYamlBaukastenV2, "1234567", "spec.components[name=mega-frontend].properties.tag"},
	{inputYamlBaukastenV2, "d821097", "spec.components[name=mega-backend].properties.tag"},
}

func TestFindNode(t *testing.T) {
	for _, test := range findNodeTests {
		nodes, err := unmarshal(test.inputYaml)
		checkErrorAndFail(err, t, "failed to unmarshal: '%s'", test.inputYaml)

		node, err := FindNode(nodes.Content[0], test.yamlPath)
		checkErrorAndFail(err, t, "node %s not found", test.yamlPath)

		if node.Value != test.expectedValue {
			t.Fatalf("expected value '%s' did not match '%s'", test.expectedValue, node.Value)
		}
	}
}

func TestFindNodeNotFound(t *testing.T) {
	yamlPath := "did.not.exist"
	nodes, err := unmarshal(inputYamlV1)
	checkErrorAndFail(err, t, "failed to unmarshal: '%s'", inputYamlV1)

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

// endregion

// region DeleteEnvFromApplicationset Tests
type deleteEnvTest struct {
	inputYaml, expectedYaml, env string
}

var deleteEnvTests = []deleteEnvTest{
	{applicationSetYaml, expectedAppSetFirst, "dev"},       // delete first element
	{applicationSetYaml, expectedAppSetMiddle, "qa"},       // delete middle element
	{applicationSetYaml, expectedAppSetLast, "feature-xy"}, // delete last element
}

func TestDeleteEnvFromApplicationset(t *testing.T) {
	for _, test := range deleteEnvTests {
		nodes, err := unmarshal(test.inputYaml)
		checkErrorAndFail(err, t, "failed to unmarshal %s", test.inputYaml)
		err = DeleteEnvFromApplicationset(nodes.Content[0], test.env)
		checkErrorAndFail(err, t, "failed to delete element %s from %s", test.env, test.inputYaml)
		nodesAsString, err := marshal(nodes)
		checkErrorAndFail(err, t, "failed to marshal result")
		if test.expectedYaml != nodesAsString {
			t.Fatalf("expected %s did not match %s", test.expectedYaml, nodesAsString)
		}
	}
}

func TestDeleteEnvFromApplicationsetNotFound(t *testing.T) {
	envToDelete := "invalidEnv"
	expectedError := fmt.Sprintf("could not find environment '%s'", envToDelete)

	nodes, err := unmarshal(applicationSetYaml)
	checkErrorAndFail(err, t, "failed to unmarshal %s", applicationSetYaml)
	err = DeleteEnvFromApplicationset(nodes.Content[0], envToDelete)
	if err == nil {
		t.Fatalf("expected error, but err == nil")
	}

	if expectedError != err.Error() {
		t.Fatalf("expected error did not match current error. expected: '%s', current: '%s'", expectedError, err.Error())
	}
}

// endregion

// region FindClusterWithBranch Tests
type findClusterTest struct {
	inputYaml, branch, expectedCluster string
}

var findClusterTests = []findClusterTest{
	{applicationSetYaml, "dev", "dev"},
	{applicationSetYaml, "qa", "qa"},
	{applicationSetYaml, "main", "main"},
}

func TestFindClusterWithBranch(t *testing.T) {
	for _, test := range findClusterTests {
		nodes, err := unmarshal(test.inputYaml)
		checkErrorAndFail(err, t, "failed to unmarshal %s", test.inputYaml)
		foundCluster, err := FindClusterWithBranch(nodes.Content[0], test.branch)
		checkErrorAndFail(err, t, "did not find any entry for branch: '%s'", test.branch)
		if test.expectedCluster != foundCluster {
			t.Fatalf("expected cluster '%s' did not match fond cluster '%s'", test.expectedCluster, foundCluster)
		}
	}
}

func TestFindClusterWithBranchNotFound(t *testing.T) {
	branchToFind := "didNotExist"
	expectedError := fmt.Sprintf("no branch with name '%s' found", branchToFind)
	nodes, err := unmarshal(applicationSetYaml)
	checkErrorAndFail(err, t, "failed to unmarshal %s", applicationSetYaml)
	foundCluster, err := FindClusterWithBranch(nodes.Content[0], branchToFind)
	if foundCluster != "" {
		t.Fatalf("expected found cluster to be empty, was '%s'", foundCluster)
	}

	if err == nil {
		t.Fatal("expected error not to be nil")
	}

	if err.Error() != expectedError {
		t.Fatalf("expected error '%s' did not match current error '%s'", expectedError, err.Error())
	}
}

// endregion

func checkErrorAndFail(err error, t *testing.T, msg string, args ...any) {
	if err != nil {
		t.Errorf(msg, args...)
		t.Fatal(err.Error())
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

	return b.String(), nil
}
