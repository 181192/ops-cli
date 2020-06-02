package update

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/181192/ops-cli/pkg/download"
	"github.com/181192/ops-cli/pkg/util"
	"github.com/181192/ops-cli/pkg/util/stringutils"
	"github.com/181192/ops-cli/pkg/util/version"
	"github.com/hashicorp/go-getter"
	semver "github.com/hashicorp/go-version"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// updateCmd represents the download command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update ops-cli to latest version",
	Long:  `Update ops-cli to latest version`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return newOpsCliRelease().Update()
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
		"windows/amd64": "ops_cli_windows_amd64",
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

// Update downloads a github release if its not present in the local config folder
func (release *opsCliRelease) Update() error {

	current, err := semver.NewVersion(version.Version)
	logger.Debugf("Current version of ops-cli %s", current)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	latestURL := "https://github.com/" + release.Account + "/" + release.Name + "/releases/latest"
	logger.Debugf("Attemting to fetch latest tag from %s", latestURL)
	resp, err := client.Get(latestURL)
	if err != nil {
		logger.Warnf("Failed to get latest stable version of ops-cli %s %s", err, resp.Status)
	}
	defer resp.Body.Close()

	location, err := resp.Location()
	if err != nil {
		logger.Warnf("Failed to get latest stable version of ops-cli %s %s", err, resp.Status)
	}

	latestTag := stringutils.After(location.Path, "tag/")
	release.Version = latestTag

	latest, err := semver.NewSemver(latestTag)
	logger.Debugf("Latest release version of ops-cli %s", latest)

	if !current.LessThan(latest) {
		logger.Infof("No updates available, using %s and latest version is %s", current.String(), latest.String())
		return nil
	}

	url := "https://github.com/" + release.Account + "/" + release.Name + "/releases/download/" + release.Version + "/" + release.ArtifactName

	progress := getter.WithProgress(download.DefaultProgressBar)

	logger.Infof("Attempting to download %s, version %s, to %q from %s", release.Name, release.Version, release.LocalFileName, url)

	tmpDir, err := ioutil.TempDir("", "ops-cli")
	if err != nil {
		return fmt.Errorf("%s\nFailed to create temp directory", err)
	}
	defer os.RemoveAll(tmpDir)

	err = getter.GetAny(tmpDir, url, progress)
	if err != nil {
		return fmt.Errorf("%s\nFailed to to download external binaries", err)
	}

	err = os.Rename(tmpDir+string(os.PathSeparator)+release.ArtifactName, release.LocalFileName)
	if err != nil {
		return fmt.Errorf("%s\nFailed to move binaries", err)
	}

	err = os.Chmod(release.LocalFileName, 0775)
	if err != nil {
		return fmt.Errorf("%s\nFailed chmod", err)
	}

	return nil
}
