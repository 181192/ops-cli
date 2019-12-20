package enable

import (
	"fmt"

	"github.com/spf13/cobra"
)

// enableCmd represents the enable command
var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable resource(s)",
	Long:  `Enable resource(s).`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("enable called")
	},
}

// Command will create the `enable` commands
func Command() *cobra.Command {

	return enableCmd
}
