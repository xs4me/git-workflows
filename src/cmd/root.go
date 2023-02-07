/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"gepaplexx/git-workflows/model"
	"gepaplexx/git-workflows/utils"
	"os"

	"github.com/spf13/cobra"
)

var Version string
var Config model.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "git-workflows",
	Short: "handle git operations for argo workflows",
	Long: `cli application to handle git operations in argo workflows. 
	
Subcommands reflect the stages in the workflow. `,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	// global flags
	rootCmd.PersistentFlags().BoolVar(&Config.Development, "dev", false, "enable development mode")
	rootCmd.PersistentFlags().StringVarP(&Config.BaseDir, "path", "p", "/mnt/out", "base directory for all operations")
	rootCmd.PersistentFlags().BoolVarP(&Config.LegacyBehavior, "legacy", "l", false, "use legacy behavior")
	rootCmd.PersistentFlags().StringVar(&Config.Username, "commit-user", "argo-ci", "username for git operations")
	rootCmd.PersistentFlags().StringVar(&Config.Email, "commit-email", "argo-ci@geaprdec.com", "email for git operations")

	rootCmd.PersistentFlags().StringVarP(&Config.GitUrl, "url", "u", "", "git url for the repository")
	rootCmd.PersistentFlags().StringVar(&Config.Reponame, "name", "", "name of the repository")
	rootCmd.PersistentFlags().StringVarP(&Config.Branch, "branch", "b", "main", "branch to checkout")

	checkoutCmd.PersistentFlags().BoolVar(&Config.Extract, "extract", false, "Extract Information about the last commiter")

}

func developmentMode(c *model.Config) {
	c.BaseDir = "./tmp/"
	c.GitUrl = "git@github.com:gepaplexx-demos/demo-microservice.git"
	c.Reponame = "demo-microservice"
	c.Branch = "test"
	c.Extract = true
	err := os.RemoveAll(c.BaseDir)
	utils.CheckIfError(err)
}
