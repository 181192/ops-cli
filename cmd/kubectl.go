package cmd

import (
	"os"

	"github.com/spf13/cobra"

	kcmd "k8s.io/kubernetes/pkg/kubectl/cmd"

	// Import to initialize client auth plugins.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// kubectlCmd represents the kubectl command
var kubectlCmd = &cobra.Command{
	Use:   "kubectl",
	Short: "kubectl controls the Kubernetes cluster manager",
	Long:  `kubectl controls the Kubernetes cluster manager.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		os.Args = append(os.Args[:1], os.Args[2:]...)
		command := kcmd.NewDefaultKubectlCommand()
		return command.Execute()
	},
}

func init() {
	rootCmd.AddCommand(kubectlCmd)
}
