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
package cmd

import (
	"errors"

	"github.com/donovanmods/7dtd-modtools/modlet"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var PackCmd = &cobra.Command{
	Use:     "pack <mod-dir>",
	Args:    cobra.ExactArgs(1),
	Short:   "Creates (packs) a modlet from the given mod directory",
	Long:    `Creates (packs) a new modlet (mod + template) from an existing Mod directory`,
	GroupID: "cmd",
	Run: func(cmd *cobra.Command, mods []string) {
		if len(mods) == 0 || mods[0] == "" {
			cobra.CheckErr(errors.New("mod directory is required"))
		}
		mdir := mods[0]

		outdir := viper.GetString("outdir")
		if outdir == "" {
			cobra.CheckErr("outdir is required")
		}

		cobra.CheckErr(modlet.Pack(mdir, outdir))
	},
}

func init() {
	PackCmd.Flags().StringP("compress", "z", "", "Compress the final modlet file")
	cobra.CheckErr(viper.BindPFlag("compress", PackCmd.Flags().Lookup("compress")))
}
