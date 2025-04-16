/*
Copyright © 2025 Donovan C. Young <dyoung522@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
*/
package new

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/donovanmods/7dtd-gamedata/modinfo"
	"github.com/donovanmods/7dtd-modtools/lib/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ModletFolder struct {
	Name string
	Path string
}

func (M *ModletFolder) Config() string {
	return filepath.Join(M.Path, "Config")
}

func Run(name string, dir string) error {
	modlet := ModletFolder{
		Name: name,
		Path: filepath.Join(dir, name),
	}

	if !viper.GetBool("force") {
		if _, err := os.Stat(modlet.Path); !os.IsNotExist(err) {
			return fmt.Errorf("modlet directory %q already exists -- refusing to overwrite", modlet.Path)
		}
	}

	cobra.CheckErr(os.MkdirAll(modlet.Config(), 0755))
	cobra.CheckErr(os.WriteFile(filepath.Join(modlet.Config(), ".keep"), []byte{}, 0644))

	logger.Info("Created %s directories\n", modlet.Path)

	xml, err := modinfo.NewModInfo(modlet.Name).XML()
	if err != nil {
		return err
	}
	cobra.CheckErr(os.WriteFile(filepath.Join(modlet.Path, "ModInfo.xml"), []byte(xml), 0644))

	logger.Debug("Wrote ModInfo.xml")

	readme := []byte(fmt.Sprintf("# %s\n\nThis is the README for a new modlet created by the 7 Days Modlet Tools (7dtd-modtools).\n", modlet.Name))
	cobra.CheckErr(os.WriteFile(filepath.Join(modlet.Path, "README.md"), readme, 0644))

	logger.Debug("Wrote README.md")

	logger.Info("Created new modlet: %s", modlet.Name)

	return nil
}
