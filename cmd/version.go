package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	// RootCmd.SetVersionTemplate(version())
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long:  `All software has versions. This is ours.`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(longVersion())
	},
}

func longVersion() string {
	return fmt.Sprintf("7DMT v%s - Donovan C. Young\n\n%s", RootCmd.Version, RootCmd.Short)
}
