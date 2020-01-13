package ssh

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	kubeconfig    string
	configContext string
)

// sshCmd represents the create command
var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "ssh resource(s)",
	Long:  `ssh resource(s).`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ssh called")
	},
}

// Command will ssh the `create` commands
func Command() *cobra.Command {
	sshCmd.AddCommand(nodeSSHCmd())

	sshCmd.PersistentFlags().StringVarP(&kubeconfig, "kubeconfig", "c", "", "Kubernetes configuration file")
	sshCmd.PersistentFlags().StringVar(&configContext, "context", "", "The name of the kubeconfig context to use")

	return sshCmd
}
