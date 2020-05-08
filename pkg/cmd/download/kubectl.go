package download

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/181192/ops-cli/pkg/download"
	"github.com/181192/ops-cli/pkg/util"
	"github.com/hashicorp/go-getter"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// https://storage.googleapis.com/kubernetes-release/release/v1.16.0/bin/linux/amd64/kubectl
// https://storage.googleapis.com/kubernetes-release/release/v1.16.0/bin/darwin/amd64/kubectl
// https://storage.googleapis.com/kubernetes-release/release/v1.16.0/bin/windows/amd64/kubectl.exe

var kubectlBinary = util.GetConfigDirectory() + "/bin/kubectl"

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
		return util.RequireFile(kubectlBinary)
	},
	Run: func(cmd *cobra.Command, args []string) {
		util.ExecuteCmd(cmd, kubectlBinary, args)
	},
}

func (release *kubectlRelease) setDownloadURL() *kubectlRelease {
	if release.Version == "latest" {
		kubectlVersion := "v1.16.0"

		resp, err := http.Get("https://storage.googleapis.com/kubernetes-release/release/stable.txt")
		if err != nil {
			logger.Warnf("Failed to get latest stable version of kubectl %s", err)
			release.Version = kubectlVersion
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			logger.Warnf("Failed to read latest version, default to %s", kubectlVersion)
			release.Version = kubectlVersion
		}
		release.Version = strings.TrimSpace(string(body))
	}
	release.URL = "https://storage.googleapis.com/kubernetes-release/release/" + release.Version + "/" + release.ArtifactName
	return release
}

func (release *kubectlRelease) DownloadIfNotExists() error {
	if _, err := os.Stat(release.LocalFileName); os.IsNotExist(err) {
		progress := getter.WithProgress(download.DefaultProgressBar)

		logger.Infof("Attempting to download %s, version %s, to %q\n", release.Name, release.Version, release.LocalFileName)

		err := getter.GetFile(release.LocalFileName, release.URL, progress)
		if err != nil {
			return fmt.Errorf("%s\nFailed to to download external binaries", err)
		}

		err = os.Chmod(release.LocalFileName, 0775)
		if err != nil {
			return fmt.Errorf("%s\nFailed chmod", err)
		}
	} else {
		logger.Infof("%s already exists at %s\n", release.Name, release.LocalFileName)
	}
	return nil
}
