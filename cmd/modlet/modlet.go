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
	sub2 "github.com/donovanmods/7dtd-modtools/cmd/modlet/new"
	sub1 "github.com/donovanmods/7dtd-modtools/cmd/modlet/pack"
	sub3 "github.com/donovanmods/7dtd-modtools/cmd/modlet/unpack"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ModletCmd = &cobra.Command{
	Use:     "modlet",
	Short:   "modlet commands",
	GroupID: "cmd",
}

func init() {
	ModletCmd.AddGroup(&cobra.Group{ID: "cmd", Title: "Commands"})

	ModletCmd.AddCommand(sub1.PackCmd)
	ModletCmd.AddCommand(sub2.NewCmd)
	ModletCmd.AddCommand(sub3.UnpackCmd)

	ModletCmd.PersistentFlags().StringP("gamedir", "G", "", "Directory where the game's Config XML files live")
	cobra.CheckErr(viper.BindPFlag("gamedir", ModletCmd.PersistentFlags().Lookup("gamedir")))
	_ = ModletCmd.MarkFlagRequired("gamedir")

	ModletCmd.PersistentFlags().StringP("output", "o", "", "Output directory and/or file")
	cobra.CheckErr(viper.BindPFlag("output", ModletCmd.PersistentFlags().Lookup("output")))

	ModletCmd.PersistentFlags().BoolP("force", "F", false, "Force overwrite of existing files")
	cobra.CheckErr(viper.BindPFlag("force", ModletCmd.PersistentFlags().Lookup("force")))
}
