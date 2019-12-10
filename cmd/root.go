package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var cfgFolder = getHome() + "/.ops"

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ops",
	Short: "ops-cli is a wrapper for devops tools",
	Long:  `A wrapper for multiple devops tools...`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ops/ops.yaml)")
	RootCmd.PersistentFlags().MarkHidden("config")
}

func getHome() string {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return home
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		if _, err := os.Stat(cfgFolder); os.IsNotExist(err) {
			os.Mkdir(cfgFolder, os.ModeDir)
		}

		viper.AddConfigPath(cfgFolder)
		viper.SetConfigName("ops")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
