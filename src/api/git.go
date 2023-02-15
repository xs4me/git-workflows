package api

import (
	"fmt"
	"gepaplexx/git-workflows/logger"
	"gepaplexx/git-workflows/model"
	"gepaplexx/git-workflows/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
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
	return repo
}

func DeployFromTo(c *model.Config, repo *git.Repository) {
	logger.Info("Deploying from %s to %s", c.FromBranch, c.ToBranch)

	fromIndex := utils.IndexOf(c.FromBranch, c.Stages)
	toIndex := utils.IndexOf(c.ToBranch, c.Stages)

	mergeable(c, fromIndex, toIndex)

	for fromIndex < toIndex {
		fromBranch := c.Stages[fromIndex]
		toBranch := c.Stages[fromIndex+1]
		merge(c, repo, fromBranch, toBranch)
		fromIndex++
	}
}

// todo: über go implementierung lösen.
func merge(c *model.Config, repo *git.Repository, fromBranch string, toBranch string) {
	wt := checkout(repo, toBranch, false)
	cmd := exec.Command("git", "config", "--global", "user.email", c.Email)
	_ = execute(cmd)
	cmd = exec.Command("git", "config", "--global", "user.name", c.Username)
	_ = execute(cmd)
	cmd = exec.Command("git", "merge", fromBranch)
	_ = execute(cmd)
	commitAndPush(c, wt, repo, fmt.Sprintf("Merge from %s to %s", fromBranch, toBranch))
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

func ExtractGitInformation(c *model.Config) {
	logger.Info("Extracting git information")
	repo, err := git.PlainOpen(c.ApplicationClonePath())
	utils.CheckIfError(err)

	r, _ := repo.Head()
	commit, err := repo.CommitObject(r.Hash())
	utils.CheckIfError(err)

	writeCommitInformation(c, "hash", commit.Hash.String()[0:7])
	writeCommitInformation(c, "user", commit.Author.Name)
	writeCommitInformation(c, "email", commit.Author.Email)
}

func writeCommitInformation(c *model.Config, typ string, content string) {
	logger.Debug("Writing commit information to file: %s", typ)
	f, err := os.Create(c.BaseDir + "/commit_" + typ)
	utils.CheckIfError(err)
	defer func(f *os.File) {
		err := f.Close()
		utils.CheckIfError(err)
	}(f)
	_, err = f.WriteString(content)
	utils.CheckIfError(err)
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
