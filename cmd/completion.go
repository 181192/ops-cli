package cmd

import (
	"fmt"
	"io"
	"strings"

	cmdUtil "github.com/181192/ops-cli/pkg/util"
	"github.com/spf13/cobra"
)

var completionShells = map[string]func(out io.Writer) error{
	"bash":       rootCmd.GenBashCompletion,
	"zsh":        rootCmd.GenZshCompletion,
	"powershell": rootCmd.GenPowerShellCompletion,
}

func completionShellsArray() []string {
	shells := make([]string, 0, len(completionShells))
	for shell := range completionShells {
		shells = append(shells, shell)
	}

	return shells
}

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion [" + strings.Join(completionShellsArray(), "|") + "]",
	Short: "Generates shell completion scripts for the specified shell",
	Long:  "Generates shell completion scripts for the specified shell.",
	Example: `To load completion run

. <(bitbucket completion bash)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(bitbucket completion bash)
`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: completionShellsArray(),
	Run: func(cmd *cobra.Command, args []string) {
		if err := runCompletion(rootCmd.OutOrStdout(), cmd, args); err != nil {
			fmt.Fprintln(cmd.OutOrStderr(), err)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}

func runCompletion(out io.Writer, cmd *cobra.Command, args []string) error {
	shells := strings.Join(completionShellsArray(), ", ")
	run, found := completionShells[args[0]]
	if !found {
		return cmdUtil.UsageErrorf(cmd, "Unsupported shell type %q. Valid shells are one of: %s", args[0], shells)
	}

	return run(out)
}
