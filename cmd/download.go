package cmd

import (
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
		newHelmfileRelease().DownloadIfNotExists()
		newHelmRelease().DownloadIfNotExists()
		newKubectlRelease().DownloadIfNotExists()
		newTerraformRelease().DownloadIfNotExists()
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
