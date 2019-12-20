package list

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List resource(s)",
	Long:  `List resource(s).`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list called")
	},
}

// Command will create the `list` commands
func Command() *cobra.Command {

	return listCmd
}
