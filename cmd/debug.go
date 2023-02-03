package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "prints current configuration of the application",
	//Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		printConfig()
	},
}

func printConfig() {

	fmt.Printf("Config: %+v", Config)
}
