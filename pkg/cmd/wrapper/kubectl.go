package wrapper

import (
	"github.com/181192/ops-cli/pkg/util"
	"github.com/181192/ops-cli/pkg/wrapper"
	"github.com/spf13/cobra"
)

// kubectlCmd represents the kubectl command
var kubectlCmd = &cobra.Command{
	Use:   "kubectl",
	Short: "kubectl controls the Kubernetes cluster manager",
	Long:  `kubectl controls the Kubernetes cluster manager.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return util.RequireFile(wrapper.KubectlBinary)
	},
	Run: func(cmd *cobra.Command, args []string) {
		util.ExecuteCmd(cmd, wrapper.KubectlBinary, args)
	},
}
