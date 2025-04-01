package cmd

import (
	"github.com/donovanmods/7dtd-modtools/lib/builder"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var flags struct {
	gamedir string
	outdir  string
}

var ModletCmd = &cobra.Command{
	Use:     "modlet <modlet-templates>",
	Args:    cobra.MinimumNArgs(1),
	Short:   "Build modlet from modlet template",
	Long:    `Builds a modlet from a modlet template file`,
	GroupID: "cmd",
	Run:     execute,
}

func execute(cmd *cobra.Command, templates []string) {
	cobra.CheckErr(builder.BuildModlet(templates, flags.gamedir, flags.outdir))
}

func init() {
	ModletCmd.Flags().StringVarP(&flags.gamedir, "gamedir", "g", "", "Directory where the game's Config XML files live")
	cobra.CheckErr(viper.BindPFlag("gamedir", ModletCmd.Flags().Lookup("gamedir")))
	_ = ModletCmd.MarkFlagRequired("gamedir")

	ModletCmd.Flags().StringVarP(&flags.outdir, "outdir", "o", ".", "Output directory")
	cobra.CheckErr(viper.BindPFlag("outdir", ModletCmd.Flags().Lookup("outdir")))

	// ModletCmd.Flags().StringVarP(&flags.name, "name", "n", "", "Name of the output modlet")
	// _ = ModletCmd.MarkFlagRequired("gamedir")
}
