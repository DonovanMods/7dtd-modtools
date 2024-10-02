/*
Copyright © 2024 Donovan C. Young <dyoung522@gmail.com>
*/
package cmd

import (
	"os"
	"path/filepath"

	sub1 "github.com/donovanmods/7dmt/cmd/make"
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:     "7dmt <command> [flags]",
	Version: "0.1.0",
	Short:   "tools used to create, modify, install, and validate 7 Days to Die Mods",
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
	RootCmd.AddCommand(sub1.MakeCmd)
	// RootCmd.AddCommand(sub2.validateCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Set the configFile if not set by flags
	if configFile == "" {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		configFile = filepath.Join(home, ".config", "7dmt", "config")
	}

	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv() // read in environment variables that match

	err := viper.ReadInConfig()
	cobra.CheckErr(err)
}
