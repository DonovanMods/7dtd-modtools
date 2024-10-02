package cmd

import (
	sub1 "github.com/donovanmods/7dmt/cmd/make/lessgrind"
	"github.com/spf13/cobra"
)

// ListCmd represents the list command
var MakeCmd = &cobra.Command{
	Use:   "make",
	Short: "Various make commands",
}

func init() {
	MakeCmd.AddCommand(sub1.LessGrindCmd)
}
