package cmd

import (
	"gepaplexx/git-workflows/api"
	"gepaplexx/git-workflows/model"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(descriptorCmd)
}

var descriptorCmd = &cobra.Command{
	Use:   "descriptor",
	Short: "checkout workflow descriptor file only",
	Run: func(cmd *cobra.Command, args []string) {
		descriptor(&Config)
	},
}

func descriptor(c *model.Config) {
	if c.Development {
		developmentMode(c)
		c.DefaultDescriptorLocation = "templates/default-descriptor.json"
	}

	descriptorPreRequisites(c)

	api.GetWorkflowDescriptor(c)
}

func descriptorPreRequisites(c *model.Config) {
	prerequisites(c)
}
