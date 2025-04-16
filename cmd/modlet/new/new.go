package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/donovanmods/7dtd-gamedata/modinfo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var flags struct {
	dir string
}

// ListCmd represents the list command
var NewCmd = &cobra.Command{
	Use:     "new [flags] <modlet-name>",
	Short:   "Create a new modlet in the location specified or in the current directory",
	Args:    cobra.ExactArgs(1),
	GroupID: "cmd",
	Run:     execute,
}

func execute(cmd *cobra.Command, name []string) {
	verbosity := viper.GetInt("verbosity")

	if len(name) == 0 {
		fmt.Println(errors.New("please provide a Modlet name to use on the command line"))
		fmt.Printf("Usage: %s\n\n", cmd.Use)
		os.Exit(1)
	}

	modletName := name[0]
	modletBase := filepath.Join(flags.dir, name[0])
	configDir := filepath.Join(modletBase, "Config")

	if !cmd.Flag("force").Changed {
		if _, err := os.Stat(modletBase); !os.IsNotExist(err) {
			cobra.CheckErr(fmt.Errorf("modlet directory %q already exists -- refusing to overwrite", modletBase))
		}
	}

	cobra.CheckErr(os.MkdirAll(configDir, 0755))
	cobra.CheckErr(os.WriteFile(filepath.Join(configDir, ".keep"), []byte{}, 0644))

	if verbosity > 1 {
		fmt.Printf("Created %s directories\n", modletBase)
	}

	xml, err := modinfo.NewModInfo(modletName).XML()
	if err != nil {
		cobra.CheckErr(err)
	}
	cobra.CheckErr(os.WriteFile(filepath.Join(modletBase, "ModInfo.xml"), []byte(xml), 0644))

	if verbosity > 1 {
		fmt.Println("Wrote ModInfo.xml")
	}

	readme := []byte(fmt.Sprintf("# %s\n\nThis is the README for a new modlet created by the 7 Days Modlet Tools (7dtd-modtools).\n", modletName))
	cobra.CheckErr(os.WriteFile(filepath.Join(modletBase, "README.md"), readme, 0644))

	if verbosity > 1 {
		fmt.Println("Wrote README.md")
	}

	if verbosity > 0 {
		fmt.Println("Created new modlet:", modletBase)
	}
}
