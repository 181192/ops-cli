package download

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/181192/ops-cli/pkg/download"
	"github.com/181192/ops-cli/pkg/util/stringutils"
	"github.com/181192/ops-cli/pkg/wrapper"
	"github.com/hashicorp/go-getter"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

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
			LocalFileName: wrapper.HelmBinary,
		},
	}

	release.SetVersion()
	release.SetArtifactName(artifactsNames)
	release.setDownloadURL()

	return release
}

// helmCmd represents the helm command
var helmCmd = &cobra.Command{
	Use:   "helm",
	Short: "Downloads helm",
	Long:  `Downloads helm.`,
	Run: func(cmd *cobra.Command, args []string) {
		newHelmRelease().DownloadIfNotExists()
	},
}

// https://get.helm.sh/helm-v3.0.1-linux-amd64.tar.gz
// https://get.helm.sh/helm-v3.0.1-darwin-amd64.tar.gz
// https://get.helm.sh/helm-v3.0.1-windows-amd64.zip

func (release *helmRelease) setDownloadURL() *helmRelease {
	if release.Version == "latest" {
		release.Version = "v3.2.1"

		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}}

		resp, err := client.Get("https://github.com/helm/helm/releases/latest")
		if err != nil {
			logger.Warnf("Failed to get latest stable version of helm %s %s", err, resp.Status)
		}
		defer resp.Body.Close()

		location, err := resp.Location()
		if err != nil {
			logger.Warnf("Failed to get latest stable version of helm %s %s", err, resp.Status)
		}

		if location.Path != "" {
			release.Version = stringutils.After(location.Path, "tag/")
		}
	}
	release.URL = "https://get.helm.sh/helm-" + release.Version + "-" + release.ArtifactName
	return release
}

func (release *helmRelease) DownloadIfNotExists() error {
	if _, err := os.Stat(release.LocalFileName); os.IsNotExist(err) {
		progress := getter.WithProgress(download.DefaultProgressBar)

		logger.Infof("Attempting to download %s, version %s, to %q\n", release.Name, release.Version, release.LocalFileName)

		tmpDir, err := ioutil.TempDir("", "ops-cli")
		if err != nil {
			logger.Fatalf("%s\nFailed to create temp directory", err)
		}
		defer os.RemoveAll(tmpDir)

		err = getter.GetAny(tmpDir, release.URL, progress)
		if err != nil {
			logger.Fatalf("%s\nFailed to to download external binaries", err)
		}

		helmDirName := tmpDir + string(os.PathSeparator) + stringutils.Before(release.ArtifactName, ".") + string(os.PathSeparator) + "helm"
		logger.Debugf("Trying to move %s to %s", helmDirName, release.LocalFileName)
		err = os.Rename(helmDirName, release.LocalFileName)
		if err != nil {
			logger.Fatalf("%s\nFailed to move binaries", err)
		}

		err = os.Chmod(release.LocalFileName, 0775)
		if err != nil {
			logger.Fatalf("%s\nFailed chmod", err)
		}
	} else {
		logger.Infof("%s already exists at %s\n", release.Name, release.LocalFileName)
	}
	return nil
}
