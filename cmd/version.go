package cmd

import (
	"fmt"

	"github.com/181192/ops-cli/pkg/util/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Output the version",
	Long:  "Output the version.",
	Run: func(cmd *cobra.Command, args []string) {

		short, _ := cmd.Flags().GetBool("short")

		if short {
			fmt.Fprintln(cmd.OutOrStdout(), version.Version)
		} else {
			fmt.Fprintln(cmd.OutOrStdout(), "Version: "+version.Version+"\nGit commit: "+version.GitCommit)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolP("short", "s", false, "Only print version number")
}
