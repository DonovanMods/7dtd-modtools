package cmd

import (
	sub1 "github.com/donovanmods/7dmt/cmd/build/modlet"
	sub2 "github.com/donovanmods/7dmt/cmd/build/new"
	"github.com/spf13/cobra"
)

// ListCmd represents the list command
var BuildCmd = &cobra.Command{
	Use:     "build",
	Short:   "build commands",
	GroupID: "cmd",
}

func init() {
	BuildCmd.AddGroup(&cobra.Group{ID: "cmd", Title: "Commands"})

	BuildCmd.AddCommand(sub1.ModletCmd)
	BuildCmd.AddCommand(sub2.NewCmd)
}
