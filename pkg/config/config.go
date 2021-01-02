package config

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// CfgFile config file path
var CfgFile string

// CfgName config name
var CfgName = "ops"

// InitConfig reads in config file and ENV variables if set.
func InitConfig() {

	if CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(CfgFile)
		viper.ReadInConfig()
	} else {

		cfgDir := GetConfigDirectory()

		if _, err := os.Stat(cfgDir); os.IsNotExist(err) {
			os.Mkdir(GetConfigDirectory(), os.ModePerm)
		}

		viper.SetConfigName(CfgName)

		// Read from global config
		global := viper.New()
		global.SetConfigName(CfgName)
		global.AddConfigPath(cfgDir)
		global.ReadInConfig()

		viper.AddConfigPath(".")
		viper.MergeInConfig()

		viper.MergeConfigMap(global.AllSettings())
	}

	// Read in environment variables that match
	viper.AutomaticEnv()
}

// GetConfigDirectory returns the config directory for the executing user
func GetConfigDirectory() string {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return home + string(os.PathSeparator) + "." + CfgName
}
