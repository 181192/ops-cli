package wrapper

import (
	"github.com/spf13/cobra"
)

// wrapperCmd represents the download command
var wrapperCmd = &cobra.Command{
	Use:     "wrapper",
	Aliases: []string{"w"},
	Short:   "",
	Long:    ``,
}

// Command will create the `wrapper` commands
func Command() *cobra.Command {

	wrapperCmd.AddCommand(helmCmd)
	wrapperCmd.AddCommand(helmfileCmd)
	wrapperCmd.AddCommand(kubectlCmd)
	wrapperCmd.AddCommand(terraformCmd)

	return wrapperCmd
}
