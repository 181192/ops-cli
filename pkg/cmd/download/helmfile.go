package download

import (
	"github.com/181192/ops-cli/pkg/download"
	"github.com/181192/ops-cli/pkg/wrapper"

	"github.com/spf13/cobra"
)

func newHelmfileRelease() *download.Release {

	var artifactsNames = map[string]string{
		"darwin/amd64":  "helmfile_darwin_amd64",
		"linux/amd64":   "helmfile_linux_amd64",
		"windows/amd64": "helmfile_windows_amd64.exe",
	}

	release := &download.Release{
		Account:       "roboll",
		Name:          "helmfile",
		LocalFileName: wrapper.HelmfileBinary,
	}

	release.SetVersion(version)
	release.SetArtifactName(artifactsNames)
	release.SetDownloadURL()

	return release
}

// helmfileCmd represents the helmfile command
var helmfileCmd = &cobra.Command{
	Use:   "helmfile",
	Short: "Downloads helmfile",
	Long:  `Downloads helmfile`,
	Run: func(cmd *cobra.Command, args []string) {
		newHelmfileRelease().DownloadIfNotExists()
	},
}
