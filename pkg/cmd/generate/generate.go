package generate

import (
	"github.com/181192/ops-cli/pkg/cmd/cmdutils"

	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate resource(s)",
	Long:  `Generate resource(s).`,
}

// Command will create the `generate` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {

	cmdutils.AddResourceCmd(flagGrouping, generateCmd, generateProfileCmd)

	return generateCmd
}
