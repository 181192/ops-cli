package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/181192/ops-cli/pkg/cmd/cmdutils"
	"github.com/181192/ops-cli/pkg/cmd/completion"
	"github.com/181192/ops-cli/pkg/cmd/config"
	"github.com/181192/ops-cli/pkg/cmd/dashboard"
	"github.com/181192/ops-cli/pkg/cmd/download"
	"github.com/181192/ops-cli/pkg/cmd/enable"
	"github.com/181192/ops-cli/pkg/cmd/generate"
	"github.com/181192/ops-cli/pkg/cmd/update"
	"github.com/181192/ops-cli/pkg/cmd/wrapper"
	cliConfig "github.com/181192/ops-cli/pkg/config"

	logger "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var loglevel string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "ops",
	Short:        "ops-cli is a wrapper for devops tools",
	Long:         `A wrapper for multiple devops tools...`,
	SilenceUsage: true,
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
	flagGrouping := cmdutils.NewGrouping()
	cobra.OnInitialize(cliConfig.InitConfig)
	rootCmd.AddCommand(dashboard.Command(flagGrouping))
	rootCmd.AddCommand(generate.Command(flagGrouping))
	rootCmd.AddCommand(enable.Command(flagGrouping))
	rootCmd.AddCommand(wrapper.Command(flagGrouping))
	rootCmd.AddCommand(download.Command(flagGrouping))
	rootCmd.AddCommand(update.Command())
	rootCmd.AddCommand(completion.Command(rootCmd))
	rootCmd.AddCommand(config.Command(flagGrouping))

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := setUpLogs(os.Stdout, loglevel); err != nil {
			return err
		}
		return nil
	}

	rootCmd.PersistentFlags().StringVar(&loglevel, "log-level", logger.InfoLevel.String(), "Log level (debug, info, warn, error, fatal, panic)")
	rootCmd.PersistentFlags().StringVar(&cliConfig.CfgFile, "config", "", "config file (default is $HOME/.ops/ops.yaml)")
	rootCmd.PersistentFlags().MarkHidden("config")

	rootCmd.SetUsageFunc(flagGrouping.Usage)

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
