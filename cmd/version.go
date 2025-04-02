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
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	// RootCmd.SetVersionTemplate(version())
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long:  `All software has versions. This is ours.`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(longVersion())
	},
}

func longVersion() string {
	return fmt.Sprintf("7dtd-modtools v%s - Donovan C. Young\n\n%s", RootCmd.Version, RootCmd.Short)
}
