package completion

import (
	"os"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Command will create the `completion` commands
func Command(rootCmd *cobra.Command) *cobra.Command {
	var bashCompletionCmd = &cobra.Command{
		Use:   "bash",
		Short: "Generates bash completion scripts",
		Long: `To load completion run

source <(ops completion bash)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
source <(ops completion bash)

If you are stuck on Bash 3 (macOS) use

source /dev/stdin <<<"$(ops completion bash)"

`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.GenBashCompletion(os.Stdout)
		},
	}
	var zshCompletionCmd = &cobra.Command{
		Use:   "zsh",
		Short: "Generates zsh completion scripts",
		Long: `To configure your zsh shell, run:

mkdir -p ~/.zsh/completion/
ops completion zsh > ~/.zsh/completion/_ops

and put the following in ~/.zshrc:

fpath=($fpath ~/.zsh/completion)

`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.GenZshCompletion(os.Stdout)
		},
	}

	cmd := &cobra.Command{
		Use:   "completion",
		Short: "Generates shell completion scripts",
		Run: func(c *cobra.Command, _ []string) {
			if err := c.Help(); err != nil {
				logger.Debugf("ignoring error %q", err.Error())
			}
		},
	}

	cmd.AddCommand(bashCompletionCmd)
	cmd.AddCommand(zshCompletionCmd)

	return cmd
}