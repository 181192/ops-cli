package cmd

import (
	"fmt"
	"io"
	"strings"

	cmdUtil "github.com/181192/ops-cli/pkg/util"
	"github.com/spf13/cobra"
)

var completionShells = map[string]func(out io.Writer, cmd *cobra.Command) error{
	"bash":       runCompletionBash,
	"zsh":        runCompletionZsh,
	"powershell": runCompletionPowerShell,
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
	Use:                   "completion",
	DisableFlagsInUseLine: true,
	Short:                 "Generates shell completion scripts for the specified shell (bash, zsh or powershell)",
	Long: `To load completion run

. <(bitbucket completion bash)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(bitbucket completion bash)
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runCompletion(rootCmd.OutOrStdout(), cmd, args); err != nil {
			fmt.Fprintln(cmd.OutOrStderr(), err)
		}
	},
	ValidArgs: completionShellsArray(),
}

func init() {
	rootCmd.AddCommand(completionCmd)
}

func runCompletion(out io.Writer, cmd *cobra.Command, args []string) error {
	shells := strings.Join(completionShellsArray(), ", ")
	if len(args) == 0 {
		return cmdUtil.UsageErrorf(cmd, "Shell not specified (%s)", shells)
	}
	if len(args) > 1 {
		return cmdUtil.UsageErrorf(cmd, "Too many arguments. Expected only the shell type of (%s)", shells)
	}
	run, found := completionShells[args[0]]
	if !found {
		return cmdUtil.UsageErrorf(cmd, "Unsupported shell type %q. Valid shells are one of: %s", args[0], shells)
	}

	return run(out, cmd.Parent())
}

func runCompletionBash(out io.Writer, rootCmd *cobra.Command) error {
	return rootCmd.GenBashCompletion(out)
}

func runCompletionZsh(out io.Writer, rootCmd *cobra.Command) error {
	return rootCmd.GenZshCompletion(out)
}

func runCompletionPowerShell(out io.Writer, rootCmd *cobra.Command) error {
	return rootCmd.GenPowerShellCompletion(out)
}
