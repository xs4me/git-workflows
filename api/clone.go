package api

import (
	"gepaplexx/git-workflows/model"
	"gepaplexx/git-workflows/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"os"
)

func Clone(c *model.Config) {
	repo, err := git.PlainClone(c.LocalPath(), false, &git.CloneOptions{
		URL:           c.GitUrl,
		Progress:      os.Stdout,
		Depth:         1,
		Tags:          0,
		SingleBranch:  false,
		ReferenceName: plumbing.NewBranchReferenceName(c.Branch),
	})

	utils.CheckIfError(err)

	if c.Extract {
		r, _ := repo.Head()
		cIter, err := repo.Log(&git.LogOptions{From: r.Hash()})
		utils.CheckIfError(err)
		lastCommit, err := cIter.Next()
		utils.CheckIfError(err)

		writeCommitInformation(c, "hash", string(lastCommit.Hash.String()[0:7]))
		writeCommitInformation(c, "user", lastCommit.Author.Name)
		writeCommitInformation(c, "email", lastCommit.Author.Email)
	}
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
