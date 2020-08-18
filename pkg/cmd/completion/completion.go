package completion

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:

$ source <(ops completion bash)

# To load completions for each session, execute once:
Linux:
  $ ops completion bash > /etc/bash_completion.d/ops
MacOS:
  $ ops completion bash > /usr/local/etc/bash_completion.d/ops

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ ops completion zsh > "${fpath[1]}/_ops"

# You will need to start a new shell for this setup to take effect.

Fish:

$ ops completion fish | source

# To load completions for each session, execute once:
$ ops completion fish > ~/.config/fish/completions/ops.fish
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Help()
			os.Exit(0)
		}

		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletion(os.Stdout)
		}
	},
}

// Command will create the `completion` commands
func Command(rootCmd *cobra.Command) *cobra.Command {
	return completionCmd
}
