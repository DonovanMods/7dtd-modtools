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
	"fmt"
	"os"
	"path/filepath"

	sub1 "github.com/donovanmods/7dtd-modtools/cmd/modlet"
	"github.com/donovanmods/7dtd-modtools/lib/logger"
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:     "7dtd-modtools <command> [flags]",
	Short:   "Tools used to create, modify, install, and validate 7 Days to Die Modlets",
	Version: "0.1.2",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		noColor, _ := cmd.Flags().GetBool("noColor")
		verbosity, _ := cmd.Flags().GetCount("verbose")

		viper.Set("color", !noColor)
		viper.Set("verbosity", verbosity)

		logger.SetLogger(verbosity)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	cc.Init(&cc.Config{
		RootCmd:  RootCmd,
		Headings: cc.HiCyan + cc.Bold + cc.Underline,
		Commands: cc.HiYellow + cc.Bold,
		Example:  cc.Italic,
		ExecName: cc.Bold,
		Flags:    cc.Bold,
	})

	RootCmd.SetVersionTemplate(version())
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file")
	RootCmd.PersistentFlags().CountP("verbose", "v", "verbose output (may be repeated)")
	RootCmd.PersistentFlags().Bool("dryrun", false, "run without performing any persistent operations")
	RootCmd.PersistentFlags().Bool("color", true, "colorize output")
	RootCmd.PersistentFlags().Bool("no-color", false, "do not output ANSI color codes")

	cobra.OnInitialize(initConfig)

	_ = viper.BindPFlag("dryrun", RootCmd.PersistentFlags().Lookup("dryrun"))

	// Add subcommands

	cmdGroup := cobra.Group{ID: "cmd", Title: "Commands"}

	RootCmd.AddGroup(&cmdGroup)
	RootCmd.AddCommand(sub1.ModletCmd)
}

func version() string {
	return fmt.Sprintln(RootCmd.Version)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Set the configFile if not set by flags
	if configFile == "" {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		configFile = filepath.Join(home, ".7dtd-modtools")
	}

	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); !errors.Is(err, os.ErrNotExist) {
		fmt.Fprintf(os.Stderr, "Error reading config file: %s\n", err)
		os.Exit(1)
	}
}
