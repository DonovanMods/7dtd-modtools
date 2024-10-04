package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var flags struct {
	name string
}

var ModletCmd = &cobra.Command{
	Use:     "modlet <modlet-templates>",
	Args:    cobra.MinimumNArgs(1),
	Short:   "Build modlet from modlet template",
	Long:    `Builds a modlet from a modlet template file`,
	GroupID: "cmd",
	Run:     execute,
}

func execute(cmd *cobra.Command, modlets []string) {
	fmt.Printf(">>>> in execute for `build modlet`")
	fmt.Printf("with name: %s, modlets: %v\n", flags.name, modlets)
}

func init() {
	ModletCmd.Flags().StringVarP(&flags.name, "name", "n", "", "Name of the output modlet")
	_ = ModletCmd.MarkFlagRequired("name")
}
