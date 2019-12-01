package cmd

import (
	"fmt"
	"os/exec"

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
		helmfile := exec.Command("helmfile", args...)
		out, err := helmfile.Output()
		if err != nil {
			fmt.Fprintln(cmd.OutOrStderr(), err)
		}

		fmt.Fprintln(cmd.OutOrStdout(), string(out))
	},
}

func init() {
	rootCmd.AddCommand(helmfileCmd)
}
