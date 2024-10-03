package cmd

import (
	sub1 "github.com/donovanmods/7dmt/cmd/build/modlet"
	"github.com/spf13/cobra"
)

// ListCmd represents the list command
var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "build commands",
}

func init() {
	BuildCmd.AddCommand(sub1.ModletCmd)
}
