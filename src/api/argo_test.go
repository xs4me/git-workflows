package api

import (
	"fmt"
	"gepaplexx/git-workflows/model"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"os"
	"testing"
)

const unitTestTag = "unittest"

var stages = []string{"main", "dev", "qa", "prod"}

type updateAllStagesTest struct {
	name, repo, imageTagLocation string
}

var updateAllStagesTests = []updateAllStagesTest{
	{"demo-microservice", "git@github.com:gepaplexx-demos/demo-microservice", "image.tag"},
}

func TestUpdateAllStages(t *testing.T) {
	for _, test := range updateAllStagesTests {
		// extra function to use advantages of defer
		executeUpdateAllStages(test, t)
	}
}

// region HELPER
func executeUpdateAllStages(test updateAllStagesTest, t *testing.T) {
	// SETUP
	config := setupConfig(test.repo)
	config.TagLocation = test.imageTagLocation
	testdir, err := os.MkdirTemp("", test.name)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(testdir)

	repo, err := cloneGitRepoPlainClone(config, testdir)
	if err != nil {
		t.Fatal(err)
	}
	wt, _ := repo.Worktree()
	// replace current tag of main stage with whatever 'unitTestTag' holds.
	err = prepareMainStage(
		fmt.Sprintf("%s/apps/env/main/%s", wt.Filesystem.Root(), config.ImageTagFilename),
		config.TagLocation)
	if err != nil {
		t.Fatal("failed to prepare test, ", err)
	}

	// TEST
	UpdateAllStages(&config, wt, repo)

	// VALIDATE
	for _, stage := range stages[1:] {
		stageToValidatePath := fmt.Sprintf("%s/apps/env/%s/%s", wt.Filesystem.Root(), stage, config.ImageTagFilename)
		rootNode := ParseYaml(stageToValidatePath)
		tagNode, err := FindNode(rootNode.Content[0], config.TagLocation)
		if err != nil {
			t.Fatal(err.Error())
		}
		if tagNode.Value != unitTestTag {
			t.Fatalf("TAG '%s' did not match expected TAG '%s'", tagNode.Value, unitTestTag)
		}
	}
}

func prepareMainStage(imageTagFile string, imageTagLocation string) error {
	rootNode := ParseYaml(imageTagFile)
	node, err := FindNode(rootNode.Content[0], imageTagLocation)
	if err != nil {
		return err
	}
	node.Value = unitTestTag
	WriteYaml(rootNode, imageTagFile)
	return nil
}

func cloneGitRepoPlainClone(c model.Config, testdir string) (*git.Repository, error) {
	auth, err := setupAuth(c.SshConfigDir)
	if err != nil {
		return nil, err
	}

	repo, err := git.PlainClone(testdir, false, &git.CloneOptions{
		URL:  fmt.Sprintf("%s%s.git", c.GitUrl, c.InfraRepoSuffix),
		Auth: auth,
	})

	return repo, nil
}
func setupAuth(sshConfigDir string) (transport.AuthMethod, error) {
	privateKeyfile := fmt.Sprintf("%s/id_rsa", sshConfigDir)
	_, err := os.Stat(privateKeyfile)
	if err != nil {
		return nil, err
	}
	publicKeys, err := ssh.NewPublicKeysFromFile("git", privateKeyfile, "")
	if err != nil {
		return nil, err
	}

	return publicKeys, nil
}

func setupConfig(repo string) model.Config {
	config := getDefaultConfig()
	config.Development = true
	config.GitUrl = repo
	homedir, _ := os.UserHomeDir()
	config.SshConfigDir = fmt.Sprintf("%s/.ssh", homedir) // TODO überlegen wie das generischer gelöst werden kann, soll auch mit GH-Action funktionieren
	return config
}
func getDefaultConfig() model.Config {
	return model.Config{
		Development:               true,
		BaseDir:                   "/mnt/out/",
		Username:                  "argo-ci",
		Email:                     "argo-ci@gepardec.com",
		GitUrl:                    "",
		Reponame:                  "",
		Branch:                    "main",
		SshConfigDir:              "/workflow/.ssh/",
		RepoToken:                 "",
		InfraRepoSuffix:           "-ci",
		ImageTag:                  "",
		ImageTagFilename:          "values.yaml",
		Stages:                    stages,
		FromBranch:                "main",
		ToBranch:                  "",
		CommitRef:                 "no reference supplied",
		Force:                     false,
		ResourcesOnly:             false,
		Descriptor:                "workflow-descriptor.json",
		DefaultDescriptorLocation: "/workflow/default-descriptor.json",
	}
}

// endregion
