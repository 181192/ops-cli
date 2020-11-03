package config

import (
	"fmt"

	"github.com/181192/ops-cli/pkg/cmd/cmdutils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func pathConfigCmd(cmd *cmdutils.Cmd) {

	cmd.CobraCommand.Use = "path"
	cmd.CobraCommand.Short = "Show current config path"
	cmd.CobraCommand.Long = ""
	cmd.CobraCommand.Run = func(_ *cobra.Command, args []string) {
		fmt.Println(viper.ConfigFileUsed())
	}
}
