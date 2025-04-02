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
	Run: func(cmd *cobra.Command, templates []string) {
		cobra.CheckErr(builder.BuildModlets(templates, flags.gamedir, flags.outdir))
	},
}

func init() {
	ModletCmd.Flags().StringVarP(&flags.gamedir, "gamedir", "g", "", "Directory where the game's Config XML files live")
	cobra.CheckErr(viper.BindPFlag("gamedir", ModletCmd.Flags().Lookup("gamedir")))
	_ = ModletCmd.MarkFlagRequired("gamedir")

	ModletCmd.Flags().StringVarP(&flags.outdir, "outdir", "o", ".", "Output directory")
	cobra.CheckErr(viper.BindPFlag("outdir", ModletCmd.Flags().Lookup("outdir")))
}
