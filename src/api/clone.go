package api

import (
	"gepaplexx/git-workflows/logger"
	"gepaplexx/git-workflows/model"
	"gepaplexx/git-workflows/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"os"
)

func CloneAppRepo(c *model.Config) {
	_, err := git.PlainClone(c.ClonePath(), false, &git.CloneOptions{
		URL:           c.GitUrl,
		Progress:      os.Stdout,
		Depth:         1,
		Tags:          0,
		SingleBranch:  false,
		ReferenceName: plumbing.NewBranchReferenceName(c.Branch),
		// TODO: implement SSH Authentication first, then handle password/deployment keys
		//Auth:
	})

	utils.CheckIfError(err)
}

func CloneInfraRepo(c *model.Config) {
	//TODO: implement
	logger.Fatal("not implemented yet")
}

func ExtractGitInformation(c *model.Config) {
	repo, err := git.PlainOpen(c.ClonePath())
	utils.CheckIfError(err)

	r, _ := repo.Head()
	commit, err := repo.CommitObject(r.Hash())
	utils.CheckIfError(err)

	writeCommitInformation(c, "hash", commit.Hash.String()[0:7])
	writeCommitInformation(c, "user", commit.Author.Name)
	writeCommitInformation(c, "email", commit.Author.Email)
}

func writeCommitInformation(c *model.Config, typ string, content string) {

	f, err := os.Create(c.BaseDir + "commit_" + typ)
	utils.CheckIfError(err)
	defer func(f *os.File) {
		err := f.Close()
		utils.CheckIfError(err)
	}(f)
	_, err = f.WriteString(content)
	utils.CheckIfError(err)
}
