package cmd

import (
	"fmt"
	"os"

	cmdUtil "github.com/181192/ops-cli/pkg/util"

	getter "github.com/hashicorp/go-getter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var helmfileBinary = cfgFolder + "/bin/helmfile"

var helmfile = GitHubRelease{
	account:  "roboll",
	name:     "helmfile",
	filename: "helmfile_linux_amd64",
	version:  getVersion(),
}

// helmfileCmd represents the helmfile command
var helmfileCmd = &cobra.Command{
	Use:   "helmfile",
	Short: "Deploy Kubernetes Helm Charts",
	Long: `Helmfile is a declarative spec for deploying helm charts. It lets you...

	- Keep a directory of chart value files and maintain changes in version control.
	- Apply CI/CD to configuration changes.
	- Periodically sync to avoid skew in environments.`,
	Run: func(cmd *cobra.Command, args []string) {
		downloadIfNotExists(helmfile)
		cmdUtil.ExecuteCmd(cmd, helmfileBinary, args)
	},
}

func init() {
	rootCmd.AddCommand(helmfileCmd)
}

// GitHubRelease represents a github release
type GitHubRelease struct {
	account  string
	name     string
	filename string
	version  string
}

// getDownloadURL based on if version is latest or tagged
func getDownloadURL(release GitHubRelease) string {
	if release.version == "latest" {
		return "https://github.com/" + release.account + "/" + release.name + "/releases/latest/download/" + release.filename
	}
	return "https://github.com/" + release.account + "/" + release.name + "/releases/download/" + release.version + "/" + release.filename
}

func getVersion() string {
	version := viper.GetString("apps.helmfile")
	if version == "" {
		return "latest"
	}
	return version
}

func downloadIfNotExists(release GitHubRelease) {
	if _, err := os.Stat(helmfileBinary); os.IsNotExist(err) {
		progress := getter.WithProgress(cmdUtil.DefaultProgressBar)
		url := getDownloadURL(release)

		fmt.Println("Attempting to download helmfile, version: " + release.version)
		err := getter.GetFile(helmfileBinary, url, progress)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Chmod(helmfileBinary, 0775)
	}
}
