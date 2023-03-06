package cmd

import (
	"gepaplexx/git-workflows/api"
	"gepaplexx/git-workflows/model"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(checkoutCmd)
}

var checkoutCmd = &cobra.Command{
	Use:   "checkout",
	Short: "handles the checkout stage of the workflow",
	Run: func(cmd *cobra.Command, args []string) {
		checkout(&Config)
	},
}

func checkout(c *model.Config) {
	if c.Development {
		developmentMode(c)
		c.Branch = "main"
		c.Env = "main"
	}

	checkoutPreRequisites(c)

	api.CloneRepo(c, c.Branch, true)
}

func checkoutPreRequisites(c *model.Config) {
	prerequisites(c)
}
