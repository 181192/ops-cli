package download

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// TODO get cfgFolder from root
var cfgFolder = getHome() + "/.ops"

func getHome() string {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return home
}

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
func Command() *cobra.Command {

	downloadCmd.AddCommand(helmCmd)
	downloadCmd.AddCommand(helmfileCmd)
	downloadCmd.AddCommand(kubectlCmd)
	downloadCmd.AddCommand(terraformCmd)

	return downloadCmd
}
