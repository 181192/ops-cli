package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version string
var gitCommit string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Output the version of sail",
	Run: func(cmd *cobra.Command, args []string) {
		if version == "" {
			version = "unversioned"
		}

		if gitCommit == "" {
			gitCommit = "master"
		}

		short, _ := cmd.Flags().GetBool("short")

		if short {
			fmt.Fprintln(cmd.OutOrStdout(), version)
		} else {
			fmt.Fprintln(cmd.OutOrStdout(), "Version: "+version+"\nGit commit: "+gitCommit)
		}
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolP("short", "s", false, "Only print version number")
}
