package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var checkoutCmd = &cobra.Command{
	Use:   "checkout",
	Short: "handles the checkout stage of the workflow",
	//Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		checkout()
	},
}

func checkout() {
	fmt.Println("checkout git@github.com/gepaplexx/gepaplexx-demos.git")
}
