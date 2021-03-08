package config

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// Config config options
type Config struct {
	// CfgFile config file path
	CfgName string

	// CfgName config name
	CfgFile string
}

var c *Config

func init() {
	c = &Config{}
}

// InitConfig reads in config file and ENV variables if set.
func InitConfig() {
	c.InitConfig()
}

// GetConfig gets the global Config instance.
func GetConfig() *Config {
	return c
}

// SetConfigName set config name
func SetConfigName(name string) {
	c.SetConfigName(name)
}

// SetConfigName set config name
func (config *Config) SetConfigName(name string) {
	if name != "" {
		config.CfgName = name
	}
}

// SetConfigFile set config name
func SetConfigFile(file string) {
	c.SetConfigFile(file)
}

// SetConfigFile set config name
func (config *Config) SetConfigFile(file string) {
	if file != "" {
		config.CfgFile = file
	}
}

// InitConfig reads in config file and ENV variables if set.
func (config *Config) InitConfig() {
	if config.CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(config.CfgFile)
		viper.ReadInConfig()
	} else {

		cfgDir := config.GetConfigDirectory()

		if _, err := os.Stat(cfgDir); os.IsNotExist(err) {
			os.Mkdir(config.GetConfigDirectory(), os.ModePerm)
		}

		viper.SetConfigName(config.CfgName)

		// Read from global config
		global := viper.New()
		global.SetConfigName(config.CfgName)
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
	return c.GetConfigDirectory()
}

// GetConfigDirectory returns the config directory for the executing user
func (config *Config) GetConfigDirectory() string {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return home + string(os.PathSeparator) + "." + config.CfgName
}
