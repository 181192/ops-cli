package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/181192/ops-cli/pkg/cmd/completion"
	"github.com/181192/ops-cli/pkg/cmd/create"
	"github.com/181192/ops-cli/pkg/cmd/dashboard"
	"github.com/181192/ops-cli/pkg/cmd/download"
	"github.com/181192/ops-cli/pkg/cmd/enable"
	"github.com/181192/ops-cli/pkg/cmd/generate"
	"github.com/181192/ops-cli/pkg/cmd/list"

	homedir "github.com/mitchellh/go-homedir"
	logger "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var cfgFolder = getHome() + "/.ops"
var loglevel string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ops",
	Short: "ops-cli is a wrapper for devops tools",
	Long:  `A wrapper for multiple devops tools...`,
}

// NewRootCmd returns a new root cmd
func NewRootCmd() *cobra.Command {
	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.EnableCommandSorting = false
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(dashboard.Command())
	if os.Getenv("OPSCLI_EXPERIMENTAL") == "true" {
		rootCmd.AddCommand(enable.Command())
		rootCmd.AddCommand(download.Command())
		rootCmd.AddCommand(create.Command())
		rootCmd.AddCommand(generate.Command())
		rootCmd.AddCommand(list.Command())
	}
	rootCmd.AddCommand(completion.Command(rootCmd))

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := setUpLogs(os.Stdout, loglevel); err != nil {
			return err
		}
		return nil
	}

	rootCmd.PersistentFlags().StringVar(&loglevel, "log-level", logger.InfoLevel.String(), "Log level (debug, info, warn, error, fatal, panic)")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ops/ops.yaml)")
	rootCmd.PersistentFlags().MarkHidden("config")
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

func setUpLogs(out io.Writer, level string) error {
	logger.SetOutput(out)
	logger.SetFormatter(&logger.TextFormatter{
		FullTimestamp: true,
	})
	lvl, err := logger.ParseLevel(level)
	if err != nil {
		return err
	}
	logger.SetLevel(lvl)
	return nil
}
