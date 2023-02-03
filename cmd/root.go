/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"gepaplexx/git-workflows/model"
	"os"

	"github.com/spf13/cobra"
)

var Version string
var Config model.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "git-workflows",
	Version: Version,
	Short:   "handle git operations for argo workflows",
	Long: `cli application to handle git operations in argo workflows. 
	
Subcommands reflect the stages in the workflow. `,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	fmt.Println(Version)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(checkoutCmd)
	rootCmd.AddCommand(debugCmd)
	// global flags
	rootCmd.PersistentFlags().StringVarP(&Config.BaseDir, "path", "p", "/mnt/out", "base directory for all operations")
	rootCmd.PersistentFlags().BoolVarP(&Config.LegacyBehavior, "legacy", "l", false, "use legacy behavior")
	rootCmd.PersistentFlags().StringVar(&Config.Username, "commit-user", "argo-ci", "username for git operations")
	rootCmd.PersistentFlags().StringVar(&Config.Email, "commit-email", "argo-ci@geaprdec.com", "email for git operations")

}
