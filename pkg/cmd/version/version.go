package version

import (
	"fmt"

	"github.com/181192/ops-cli/pkg/cmd/cmdutils"
	"github.com/181192/ops-cli/pkg/util/version"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var shortVersion bool

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Output the version",
	Long:  "Output the version.",
	Run: func(cmd *cobra.Command, args []string) {
		if shortVersion {
			fmt.Fprintln(cmd.OutOrStdout(), version.Version)
			return
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Version: "+version.Version+"\nGit commit: "+version.GitCommit)
	},
}

// Command will create the `version` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {

	flagSetGroup := flagGrouping.New(versionCmd)
	flagSetGroup.InFlagSet("Version", func(fs *pflag.FlagSet) {
		fs.BoolVarP(&shortVersion, "short", "s", false, "Only print version number")
	})
	flagSetGroup.AddTo(versionCmd)

	return versionCmd
}
