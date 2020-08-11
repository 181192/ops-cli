package update

import (
	"github.com/181192/ops-cli/pkg/download"
	"github.com/181192/ops-cli/pkg/util"
	"github.com/spf13/cobra"
)

// updateCmd represents the download command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update ops-cli to latest version",
	Long:  `Update ops-cli to latest version`,
	Run: func(cmd *cobra.Command, args []string) {
		newOpsCliRelease().Update()
	},
}

// Command will create the `update` commands
func Command() *cobra.Command {
	return updateCmd
}

func newOpsCliRelease() *opsCliRelease {

	var artifactsNames = map[string]string{
		"darwin/amd64":  "ops_cli_darwin_amd64",
		"linux/amd64":   "ops_cli_linux_amd64",
		"linux/arm":     "ops_cli_linux_arm",
		"linux/arm64":   "ops_cli_linux_arm64",
		"windows/amd64": "ops_cli_windows_amd64.exe",
	}

	release := &opsCliRelease{
		&download.Release{
			Account:       "181192",
			Name:          "ops-cli",
			LocalFileName: util.GetExecutable(),
		},
	}

	release.SetArtifactName(artifactsNames)
	return release
}

type opsCliRelease struct {
	*download.Release
}
