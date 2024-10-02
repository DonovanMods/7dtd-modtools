package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.SetVersionTemplate(version())
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long:  `All software has versions. This is ours.`,

	Run: func(cmd *cobra.Command, args []string) {
		version := version()

		fmt.Print(version)
	},
}

func version() string {
	return fmt.Sprintf("7dmt v%s - Donovan C. Young\n", RootCmd.Version)
}
