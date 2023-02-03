package cmd

import (
	"fmt"
	"gepaplexx/git-workflows/model"
	"github.com/spf13/cobra"
)

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "prints current configuration of the application",
	//Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		printConfig(&Config)
	},
}

func printConfig(c *model.Config) {
	if c.Development {
		fmt.Println("Development mode enabled. Using local configuration.")
		developmentMode(c)
	}

	fmt.Printf("Config: %+v", Config)
}
