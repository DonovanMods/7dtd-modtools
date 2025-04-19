package cmd

import (
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

		file, err := modlet.ValidateOutputFile(viper.GetString("output"))
		if err != nil {
			cobra.CheckErr(err)
		}

		cobra.CheckErr(modlet.NewCmd(modlet.CmdArgs{
			Name:   name[0],
			Output: file,
			Force:  viper.GetBool("force"),
		}))
	},
}
