package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var LessGrindCmd = &cobra.Command{
	Use:   "lessgrind",
	Short: "Build lessgrind mod",
	Long:  `Builds the donovan-lessgrind mod by reading the game's XML files and applying specific changes`,
	Run:   execute,
}

func execute(cmd *cobra.Command, args []string) {
	fmt.Println(">>>> in execute for `make lessgrind`")
	fmt.Println(viper.Get("7dmt_gamedir"))
	fmt.Println(viper.Get("7dmt_modlets"))
}
