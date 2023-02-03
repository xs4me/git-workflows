package cmd

import (
	"fmt"
	"gepaplexx/git-workflows/api"
	"gepaplexx/git-workflows/model"
	"github.com/spf13/cobra"
	"os"
)

var checkoutCmd = &cobra.Command{
	Use:   "checkout",
	Short: "handles the checkout stage of the workflow",
	Run: func(cmd *cobra.Command, args []string) {
		checkout(&Config)
	},
}

func checkout(c *model.Config) {
	if c.Development {
		fmt.Println("Development mode enabled. Using local configuration.")
		developmentMode(c)
	}

	fmt.Printf("checkout %s\n", c.GitUrl)
	fmt.Println("Checking prerequisites...")
	prerequisites(c)

	api.Clone(c)
	api.Checkout(c)
}

func prerequisites(c *model.Config) {
	if c.GitUrl == "" {
		fmt.Println("GitUrl is empty. Please provide a valid git url.")
		os.Exit(1)
	}
	if c.Reponame == "" {
		fmt.Println("Repository name is empty. Please provide a valid repository name.")
	}
}
