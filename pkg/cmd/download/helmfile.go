package download

import (
	"github.com/181192/ops-cli/pkg/download"
	cmdUtil "github.com/181192/ops-cli/pkg/util"

	"github.com/spf13/cobra"
)

var helmfileBinary = cfgFolder + "/bin/helmfile"

func newHelmfileRelease() *download.Release {

	var artifactsNames = map[string]string{
		"darwin/amd64":  "helmfile_darwin_amd64",
		"linux/amd64":   "helmfile_linux_amd64",
		"windows/amd64": "helmfile_windows_amd64.exe",
	}

	release := &download.Release{
		Account:       "roboll",
		Name:          "helmfile",
		LocalFileName: helmfileBinary,
	}

	release.SetVersion()
	release.SetArtifactName(artifactsNames)
	release.SetDownloadURL()

	return release
}

// helmfileCmd represents the helmfile command
var helmfileCmd = &cobra.Command{
	Use:   "helmfile",
	Short: "Deploy Kubernetes Helm Charts",
	Long: `Helmfile is a declarative spec for deploying helm charts. It lets you...

	- Keep a directory of chart value files and maintain changes in version control.
	- Apply CI/CD to configuration changes.
	- Periodically sync to avoid skew in environments.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return cmdUtil.RequireFile(helmfileBinary)
	},
	Run: func(cmd *cobra.Command, args []string) {
		// TODO add --helm-binary flag as default
		cmdUtil.ExecuteCmd(cmd, helmfileBinary, args)
	},
}
