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
	sub1 "github.com/donovanmods/7dtd-modtools/cmd/build/modlet"
	sub2 "github.com/donovanmods/7dtd-modtools/cmd/build/new"
	"github.com/spf13/cobra"
)

// ListCmd represents the list command
var BuildCmd = &cobra.Command{
	Use:     "build",
	Short:   "build commands",
	GroupID: "cmd",
}

func init() {
	BuildCmd.AddGroup(&cobra.Group{ID: "cmd", Title: "Commands"})

	BuildCmd.AddCommand(sub1.ModletCmd)
	BuildCmd.AddCommand(sub2.NewCmd)
}
