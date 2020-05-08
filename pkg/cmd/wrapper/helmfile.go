package wrapper

import (
	"github.com/181192/ops-cli/pkg/util"
	"github.com/181192/ops-cli/pkg/wrapper"
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
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return util.RequireFile(wrapper.HelmfileBinary)
	},
	Run: func(cmd *cobra.Command, args []string) {
		// TODO add --helm-binary flag as default
		util.ExecuteCmd(cmd, wrapper.HelmfileBinary, args)
	},
}
