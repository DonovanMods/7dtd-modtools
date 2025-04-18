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
	"github.com/donovanmods/7dtd-modtools/modlet"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var UnpackCmd = &cobra.Command{
	Use:     "unpack <modlet-templates>",
	Args:    cobra.MinimumNArgs(1),
	Short:   "Builds (unpacks) a mod from a modlet template",
	Long:    `Builds (unpacks) a mod directory from the given modlet (mod + template)`,
	GroupID: "cmd",
	Run: func(cmd *cobra.Command, templates []string) {
		if len(templates) == 0 {
			cobra.CheckErr("no modlet templates provided")
		}

		gamedir := viper.GetString("gamedir")
		if gamedir == "" {
			cobra.CheckErr("gamedir is required")
		}

		output := viper.GetString("output")
		if output == "" {
			cobra.CheckErr("output is required")
		}

		cobra.CheckErr(modlet.Unpack(templates, gamedir, output))
	},
}
