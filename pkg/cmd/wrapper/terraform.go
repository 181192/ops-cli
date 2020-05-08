package wrapper

import (
	"github.com/181192/ops-cli/pkg/util"
	"github.com/181192/ops-cli/pkg/wrapper"
	"github.com/spf13/cobra"
)

// terraformCmd represents the terraform command
var terraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "Terraform IaC tool",
	Long:  `Terraform IaC tool.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return util.RequireFile(wrapper.TerraformBinary)
	},
	Run: func(cmd *cobra.Command, args []string) {
		util.ExecuteCmd(cmd, wrapper.TerraformBinary, args)
	},
}
