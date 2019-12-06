package cmd

import (
	gh "github.com/181192/ops-cli/pkg/github"
	cmdUtil "github.com/181192/ops-cli/pkg/util"

	"github.com/spf13/cobra"
)

var helmfile = gh.Release{
	Account:       "roboll",
	Name:          "helmfile",
	LocalFileName: cfgFolder + "/bin/helmfile",
}

var artifactsNames = map[string]string{
	"darwin/386":    "helmfile_darwin_386",
	"darwin/amd64":  "helmfile_darwin_amd64",
	"linux/386":     "helmfile_linux_386",
	"linux/amd64":   "helmfile_linux_amd64",
	"windows/386":   "helmfile_windows_386.exe",
	"windows/amd64": "helmfile_windows_amd64.exe",
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
		return cmdUtil.RequireFile(helmfile.LocalFileName)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmdUtil.ExecuteCmd(cmd, helmfile.LocalFileName, args)
	},
}

func init() {
	rootCmd.AddCommand(helmfileCmd)
}

func getHelmfileRelease() gh.Release {
	helmfile.ArtifactName = gh.GetArtifactName(artifactsNames)
	helmfile.Version = gh.GetVersion(helmfile)
	return helmfile
}
