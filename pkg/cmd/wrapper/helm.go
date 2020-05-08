package wrapper

import (
	"github.com/181192/ops-cli/pkg/util"
	"github.com/181192/ops-cli/pkg/wrapper"
	"github.com/spf13/cobra"
)

// helmCmd represents the helm command
var helmCmd = &cobra.Command{
	Use:   "helm",
	Short: "A kubernetes package manager",
	Long:  `A kubernetes package manager.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return util.RequireFile(wrapper.HelmBinary)
	},
	Run: func(cmd *cobra.Command, args []string) {
		util.ExecuteCmd(cmd, wrapper.HelmBinary, args)
	},
}
