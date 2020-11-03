package config

import (
	"os"

	"github.com/181192/ops-cli/pkg/util"
	"github.com/spf13/viper"
)

// CfgFile config file path
var CfgFile string

// InitConfig reads in config file and ENV variables if set.
func InitConfig() {
	if CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(CfgFile)
	} else {

		cfgDir := util.GetConfigDirectory()

		if _, err := os.Stat(cfgDir); os.IsNotExist(err) {
			os.Mkdir(util.GetConfigDirectory(), os.ModePerm)
		}

		viper.AddConfigPath(cfgDir)
		viper.SetConfigName("ops")
	}

	// read in environment variables that match
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	viper.ReadInConfig()
}
