package cmd

import (
	"fmt"
	"os"

	"github.com/181192/ops-cli/pkg/download"
	cmdUtil "github.com/181192/ops-cli/pkg/util"
	"github.com/hashicorp/go-getter"
	"github.com/spf13/cobra"
)

// https://storage.googleapis.com/kubernetes-release/release/v1.16.0/bin/linux/amd64/kubectl
// https://storage.googleapis.com/kubernetes-release/release/v1.16.0/bin/darwin/amd64/kubectl
// https://storage.googleapis.com/kubernetes-release/release/v1.16.0/bin/windows/amd64/kubectl.exe

var kubectlBinary = cfgFolder + "/bin/kubectl"

type kubectlRelease struct {
	*download.Release
}

func newKubectlRelease() *kubectlRelease {

	var artifactsNames = map[string]string{
		"darwin/amd64":  "bin/darwin/amd64/kubectl",
		"linux/amd64":   "bin/linux/amd64/kubectl",
		"windows/amd64": "bin/windows/amd64/kubectl.exe",
	}

	release := &kubectlRelease{
		&download.Release{
			Account:       "kubectl",
			Name:          "kubectl",
			LocalFileName: kubectlBinary,
		},
	}

	release.SetVersion()
	release.SetArtifactName(artifactsNames)
	release.setDownloadURL()

	return release
}

// kubectlCmd represents the kubectl command
var kubectlCmd = &cobra.Command{
	Use:   "kubectl",
	Short: "kubectl controls the Kubernetes cluster manager",
	Long:  `kubectl controls the Kubernetes cluster manager.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return cmdUtil.RequireFile(kubectlBinary)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmdUtil.ExecuteCmd(cmd, kubectlBinary, args)
	},
}

func init() {
	rootCmd.AddCommand(kubectlCmd)
}

func (release *kubectlRelease) setDownloadURL() *kubectlRelease {
	if release.Version == "latest" {
		release.Version = "v1.16.0"
	}
	release.URL = "https://storage.googleapis.com/kubernetes-release/release/" + release.Version + "/" + release.ArtifactName
	return release
}

func (release *kubectlRelease) DownloadIfNotExists() error {
	if _, err := os.Stat(release.LocalFileName); os.IsNotExist(err) {
		progress := getter.WithProgress(download.DefaultProgressBar)

		fmt.Printf("Attempting to download %s, version %s, to %q\n", release.Name, release.Version, release.LocalFileName)

		err := getter.GetFile(release.LocalFileName, release.URL, progress)
		if err != nil {
			return fmt.Errorf("%s\nFailed to to download external binaries", err)
		}

		err = os.Chmod(release.LocalFileName, 0775)
		if err != nil {
			return fmt.Errorf("%s\nFailed chmod", err)
		}
	} else {
		fmt.Printf("%s already exists at %s\n", release.Name, release.LocalFileName)
	}
	return nil
}
