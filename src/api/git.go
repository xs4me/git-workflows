package api

import (
	"fmt"
	"gepaplexx/git-workflows/logger"
	"gepaplexx/git-workflows/model"
	"gepaplexx/git-workflows/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"os"
	"strings"
)

func CloneRepo(c *model.Config, appRepo bool) *git.Repository {
	checkoutPrerequisites(c)
	path, url := getCorrectRepositoryInformation(c, appRepo)

	logger.Info("Cloning Repository %s, to %s", url, path)
	repo, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:           url,
		Progress:      nil,
		Depth:         1,
		Tags:          0,
		SingleBranch:  false,
		ReferenceName: plumbing.NewBranchReferenceName(c.Branch),
		Auth:          gitAuthMethod(c),
	})

	utils.CheckIfError(err)
	return repo
}

func OpenRepo(c *model.Config, appRepo bool) *git.Repository {
	path, _ := getCorrectRepositoryInformation(c, appRepo)
	repo, err := git.PlainOpen(path)
	utils.CheckIfError(err)
	return repo
}

func getCorrectRepositoryInformation(c *model.Config, appRepo bool) (path string, url string) {
	logger.Debug("Checking if application repository is requested, appRepo: %t", appRepo)
	if appRepo {
		return c.ApplicationClonePath(), c.GitUrl
	} else {
		url := fmt.Sprintf("%s%s%s", strings.TrimSuffix(c.GitUrl, ".git"), c.InfraRepoSuffix, ".git")
		return c.InfrastructureClonePath(), url
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
	f, err := os.Create(c.BaseDir + "commit_" + typ)
	utils.CheckIfError(err)
	defer func(f *os.File) {
		err := f.Close()
		utils.CheckIfError(err)
	}(f)
	_, err = f.WriteString(content)
	utils.CheckIfError(err)
}

func checkoutPrerequisites(c *model.Config) {
	logger.Debug("Checking checkout prerequisites")
	if c.GitUrl == "" {
		logger.Error("GitUrl is empty. Please provide a valid git url.")
		os.Exit(1)
	}
	if c.Reponame == "" {
		logger.Error("Repository name is empty. Please provide a valid repository name.")
		os.Exit(1)
	}
}

func checkout(c *model.Config, repo *git.Repository) *git.Worktree {
	wt, err := repo.Worktree()
	utils.CheckIfError(err)
	err = wt.Checkout(&git.CheckoutOptions{
		Create: false,
		Branch: plumbing.NewBranchReferenceName(c.Branch),
	})
	utils.CheckIfError(err)
	return wt
}
