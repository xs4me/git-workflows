package api

import (
	"gepaplexx/git-workflows/model"
	"github.com/go-git/go-git/v5"
	"os"
)

func Clone(c *model.Config) {
	_, err := git.PlainClone(c.LocalPath(), false, &git.CloneOptions{
		URL:          c.GitUrl,
		Progress:     os.Stdout,
		Depth:        1,
		Tags:         0,
		SingleBranch: false,
	}), &git.CheckoutOptions{
		//Branch: c.Branch TODO add checkout to specific branch
	}

	if err != nil {
		panic(err)
	}

}
