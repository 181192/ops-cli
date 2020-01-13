package enable

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	kubeconfig    string
	configContext string
	namespace     string
)

// enableCmd represents the enable command
var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable resource(s)",
	Long:  `Enable resource(s).`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
}

// Command will create the `enable` commands
func Command() *cobra.Command {

	enableCmd.AddCommand(enableProfileCmd())

	enableCmd.PersistentFlags().StringVarP(&kubeconfig, "kubeconfig", "c", "", "Kubernetes configuration file")
	enableCmd.PersistentFlags().StringVar(&configContext, "context", "", "The name of the kubeconfig context to use")

	return enableCmd
}
