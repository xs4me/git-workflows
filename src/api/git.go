package api

import (
	"fmt"
	"gepaplexx/git-workflows/logger"
	"gepaplexx/git-workflows/model"
	"gepaplexx/git-workflows/utils"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/qjebbs/go-jsons"
	"os"
	"os/exec"
	"strings"
	"time"
)

func CloneRepo(c *model.Config, branch string, appRepo bool) *git.Repository {
	path, url, depth := getCorrectRepositoryInformation(c, appRepo)

	logger.Info("Cloning Repository %s, to %s", url, path)
	repo, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:             url,
		Progress:        nil,
		Depth:           depth,
		Tags:            0,
		SingleBranch:    false,
		ReferenceName:   plumbing.NewBranchReferenceName(branch),
		Auth:            gitAuthMethod(c),
		InsecureSkipTLS: true,
	})
	utils.CheckIfError(err)

	err = repo.Fetch(&git.FetchOptions{
		Auth:     gitAuthMethod(c),
		Depth:    1,
		RefSpecs: []config.RefSpec{"refs/*:refs/*"},
	})
	utils.CheckIfError(err)

	_ = checkout(repo, branch, false)
	return repo
}

func GetWorkflowDescriptor(c *model.Config) {
	logger.Info("Getting Workflow Descriptor")
	repo := CloneRepo(c, c.Branch, true)

	wt, err := repo.Worktree()
	utils.CheckIfError(err)

	defDescriptor, err := os.ReadFile("templates/default-descriptor.json")
	utils.CheckIfError(err)
	actDescriptor, err := wt.Filesystem.Open(c.Descriptor)
	utils.CheckIfError(err)
	defer func(actDescriptor billy.File) {
		err := actDescriptor.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}(actDescriptor)

	mergedDescriptor, err := jsons.Merge(defDescriptor, actDescriptor)
	utils.CheckIfError(err)

	err = os.WriteFile(fmt.Sprintf("%s/%s", c.BaseDir, "workflow-descriptor.json"), mergedDescriptor, 0644)
	utils.CheckIfError(err)
}

func DeployFromTo(c *model.Config, repo *git.Repository) {
	logger.Info("Deploying from %s to %s", c.FromBranch, c.ToBranch)

	fromIndex := utils.IndexOf(c.FromBranch, c.Stages)
	toIndex := utils.IndexOf(c.ToBranch, c.Stages)
	wt, err := repo.Worktree()
	utils.CheckIfError(err)
	dir := wt.Filesystem.Root()

	mergeable(c, fromIndex, toIndex)
	cmd := exec.Command("git", "config", "--local", "user.email", c.Email)
	cmd.Dir = dir
	_ = execute(cmd)

	cmd = exec.Command("git", "config", "--local", "user.name", c.Username)
	cmd.Dir = dir
	_ = execute(cmd)

	cmd = exec.Command("git", "config", "--local", "push.autoSetupRemote", "true")
	cmd.Dir = dir
	_ = execute(cmd)

	for fromIndex < toIndex {
		fromBranch := c.Stages[fromIndex]
		toBranch := c.Stages[fromIndex+1]
		merge(wt, fromBranch, toBranch)
		push(c, repo)
		fromIndex++
	}
}

// todo: über go implementierung lösen.
func merge(wt *git.Worktree, fromBranch string, toBranch string) {
	cmd := exec.Command("git", "checkout", toBranch)
	cmd.Dir = wt.Filesystem.Root()
	_ = execute(cmd)

	cmd = exec.Command("git", "rebase", fromBranch, toBranch)
	cmd.Dir = wt.Filesystem.Root()
	_ = execute(cmd)
}

func mergeable(c *model.Config, fromIndex int, toIndex int) {
	if fromIndex == -1 || toIndex == -1 {
		logger.Fatal("Source or target stage not in configured stages. exiting")
	}

	if fromIndex > toIndex {
		logger.Fatal("Configuration doesn't allow merge from %s to %s. exiting", c.FromBranch, c.ToBranch)
	}

	if fromIndex == toIndex {
		logger.Info("Nothing to do. exiting")
		return
	}
}

func checkout(repo *git.Repository, branch string, create bool) *git.Worktree {
	wt, err := repo.Worktree()
	utils.CheckIfError(err)
	err = wt.Checkout(&git.CheckoutOptions{
		Create: create,
		Branch: plumbing.NewBranchReferenceName(branch),
	})
	utils.CheckIfError(err)
	return wt
}

func getCorrectRepositoryInformation(c *model.Config, appRepo bool) (path string, url string, depth int) {
	logger.Debug("Checking if application repository is requested, appRepo: %t", appRepo)
	if appRepo {
		return c.ApplicationClonePath(), c.GitUrl, 1
	} else {
		url := fmt.Sprintf("%s%s%s", strings.TrimSuffix(c.GitUrl, ".git"), c.InfraRepoSuffix, ".git")
		return c.InfrastructureClonePath(), url, 0
	}
}

func commitAndPush(c *model.Config, wt *git.Worktree, repo *git.Repository, message string) {
	logger.Info("Committing and pushing changes: %s", message)
	commit(c, wt, message)
	push(c, repo)
}

func commit(c *model.Config, wt *git.Worktree, message string) {
	err := wt.AddWithOptions(&git.AddOptions{
		All: true,
	})
	utils.CheckIfError(err)
	_, err = wt.Commit(message, &git.CommitOptions{
		Committer: &object.Signature{
			Name:  c.Username,
			Email: c.Email,
			When:  time.Now(),
		},
		Author: &object.Signature{
			Name:  "argo-ci",
			Email: "argo-ci@gepardec.com",
			When:  time.Now(),
		},
	})
	utils.CheckIfError(err)
}

func push(c *model.Config, repo *git.Repository) {
	if c.IsPushEnabled() {
		logger.Info("Pushing changes to remote repository")
		err := repo.Push(&git.PushOptions{
			Auth: gitAuthMethod(c),
		})
		utils.CheckIfError(err)
	} else {
		logger.Debug("Development mode is enabled. Skipping push to remote repository")
	}
}
