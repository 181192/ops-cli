package create

import (
	"fmt"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create resource(s)",
	Long:  `Create resource(s).`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")
	},
}

// Command will create the `create` commands
func Command() *cobra.Command {

	return createCmd
}
