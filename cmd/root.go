/*
Copyright © 2024 Donovan C. Young <dyoung522@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	sub1 "github.com/donovanmods/7dtd-modtools/cmd/build"
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
		verbosity, _ := cmd.Flags().GetCount("verbose")
		viper.Set("verbosity", verbosity)
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
	// var Verbose int
	var err error

	cobra.OnInitialize(initConfig)

	// RootCmd.PersistentFlags().StringVar(&configFlag, "config", "", fmt.Sprintf("config file (default is %s)", filepath.Join(configPath, "config")))
	RootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file")
	RootCmd.PersistentFlags().CountP("verbose", "v", "counted verbosity")
	RootCmd.PersistentFlags().BoolP("no-color", "N", false, "do not output ANSI color codes")

	err = viper.BindPFlag("noColor", RootCmd.PersistentFlags().Lookup("no-color"))
	if err != nil {
		panic(err)
	}

	// Add subcommands

	cmdGroup := cobra.Group{ID: "cmd", Title: "Commands"}

	RootCmd.AddGroup(&cmdGroup)
	RootCmd.AddCommand(sub1.BuildCmd)
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
