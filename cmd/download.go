package cmd

import (
	gh "github.com/181192/ops-cli/pkg/github"

	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Downloads all required binaries to use the ops-cli",
	Long: `Downloads all required binaries to use the ops-cli.

Will download the given versions in the config file if presents
`,
	Run: func(cmd *cobra.Command, args []string) {
		gh.DownloadIfNotExists(getHelmfileRelease())
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
