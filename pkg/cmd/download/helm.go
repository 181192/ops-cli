package download

import (
	"fmt"
	"os"

	"github.com/181192/ops-cli/pkg/download"
	cmdUtil "github.com/181192/ops-cli/pkg/util"
	"github.com/hashicorp/go-getter"

	"github.com/spf13/cobra"
)

var helmBinary = cfgFolder + "/bin/helm"

type helmRelease struct {
	*download.Release
}

func newHelmRelease() *helmRelease {

	var artifactsNames = map[string]string{
		"darwin/amd64":  "darwin-amd64.tar.gz",
		"linux/amd64":   "linux-amd64.tar.gz",
		"windows/amd64": "windows-amd64.zip",
	}

	release := &helmRelease{
		&download.Release{
			Account:       "helm",
			Name:          "helm",
			LocalFileName: helmBinary,
		},
	}

	release.SetVersion()
	release.SetArtifactName(artifactsNames)
	release.setDownloadURL()

	return release
}

// helmfileCmd represents the helmfile command
var helmCmd = &cobra.Command{
	Use:   "helm",
	Short: "A kubernetes package manager",
	Long:  `A kubernetes package manager.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return cmdUtil.RequireFile(helmBinary)
	},
	Run: func(cmd *cobra.Command, args []string) {
		// TODO Set chache, configuration and data paths
		cmdUtil.ExecuteCmd(cmd, helmBinary, args)
	},
}

// https://get.helm.sh/helm-v3.0.1-linux-amd64.tar.gz
// https://get.helm.sh/helm-v3.0.1-darwin-amd64.tar.gz
// https://get.helm.sh/helm-v3.0.1-windows-amd64.zip

func (release *helmRelease) setDownloadURL() *helmRelease {
	if release.Version == "latest" {
		release.Version = "v3.0.1"
	}
	release.URL = "https://get.helm.sh/helm-" + release.Version + "-" + release.ArtifactName
	return release
}

func (release *helmRelease) DownloadIfNotExists() error {
	if _, err := os.Stat(release.LocalFileName); os.IsNotExist(err) {
		progress := getter.WithProgress(download.DefaultProgressBar)

		fmt.Printf("Attempting to download %s, version %s, to %q\n", release.Name, release.Version, release.LocalFileName)
		tmpDir := cfgFolder + "/bin/.tmp"

		err := getter.GetAny(tmpDir, release.URL, progress)
		if err != nil {
			return fmt.Errorf("%s\nFailed to to download external binaries", err)
		}

		err = os.Rename(tmpDir+"/linux-amd64/helm", release.LocalFileName)
		if err != nil {
			return fmt.Errorf("%s\nFailed to move binaries", err)
		}

		err = os.Chmod(release.LocalFileName, 0775)
		if err != nil {
			return fmt.Errorf("%s\nFailed chmod", err)
		}

		err = os.RemoveAll(tmpDir)
		if err != nil {
			return fmt.Errorf("%s\nFailed delete tmp dir", err)
		}
	} else {
		fmt.Printf("%s already exists at %s\n", release.Name, release.LocalFileName)
	}
	return nil
}
