package generate

import (
	"fmt"

	"github.com/spf13/cobra"
)

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

	return generateCmd
}
