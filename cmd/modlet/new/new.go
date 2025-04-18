package cmd

import (
	"path/filepath"

	"github.com/donovanmods/7dtd-modtools/modlet"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ListCmd represents the list command
var NewCmd = &cobra.Command{
	Use:     "new [flags] <modlet-name>",
	Short:   "Create a new modlet in the location specified or in the current directory",
	Args:    cobra.ExactArgs(1),
	GroupID: "cmd",
	Run: func(cmd *cobra.Command, name []string) {
		if len(name) == 0 {
			cobra.CheckErr("please provide a Modlet name to use on the command line")
		}

		output := filepath.Clean(viper.GetString("output"))

		cobra.CheckErr(modlet.New(name[0], output))
	},
}
