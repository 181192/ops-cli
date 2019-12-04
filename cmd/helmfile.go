package cmd

import (
	cmdUtil "github.com/181192/ops-cli/cmd/util"

	"github.com/spf13/cobra"
)

// helmfileCmd represents the helmfile command
var helmfileCmd = &cobra.Command{
	Use:   "helmfile",
	Short: "Deploy Kubernetes Helm Charts",
	Long: `Helmfile is a declarative spec for deploying helm charts. It lets you...

	- Keep a directory of chart value files and maintain changes in version control.
	- Apply CI/CD to configuration changes.
	- Periodically sync to avoid skew in environments.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmdUtil.ExecuteCmd(cmd, "helmfile", args)
	},
}

func init() {
	rootCmd.AddCommand(helmfileCmd)
}
