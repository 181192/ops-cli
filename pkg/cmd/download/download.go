package download

import (
	"github.com/181192/ops-cli/pkg/cmd/cmdutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var version string

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

// Command will create the `download` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {

	downloadCmd.AddCommand(helmCmd)
	downloadCmd.AddCommand(helmfileCmd)
	downloadCmd.AddCommand(kubectlCmd)
	downloadCmd.AddCommand(terraformCmd)

	addFlagSets(flagGrouping, []*cobra.Command{helmCmd, helmfileCmd, kubectlCmd, terraformCmd})

	return downloadCmd
}

func addFlagSets(flagGrouping *cmdutils.FlagGrouping, cmds []*cobra.Command) {
	for _, cmd := range cmds {
		flagSetGroup := flagGrouping.New(cmd)
		flagSetGroup.InFlagSet("Download", func(fs *pflag.FlagSet) {
			fs.StringVarP(&version, "version", "v", "", "Version to download (git tag)")
		})
		flagSetGroup.AddTo(cmd)
	}
}
