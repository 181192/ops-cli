package config

import (
	"github.com/181192/ops-cli/pkg/cmd/cmdutils"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Config related commands",
	Long:  `Config related commands`,
}

// Command will create the `config` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {

	cmdutils.AddResourceCmd(flagGrouping, configCmd, showConfigCmd)
	return configCmd
}
