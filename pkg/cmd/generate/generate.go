package generate

import (
	"fmt"

	"github.com/spf13/cobra"
)

var clusterConfigFile string

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate resource(s)",
	Long:  `Generate resource(s).`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("generate called")
	},
}

// Command will create the `generate` commands
func Command() *cobra.Command {

	generateCmd.AddCommand(generateProfileCmd())

	generateCmd.PersistentFlags().StringVarP(&clusterConfigFile, "cluster-config-file", "f", "", "Load configuration from a file (or stdin if set to '-')")
	return generateCmd
}
