// +build !windows

package update

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/181192/ops-cli/pkg/download"
	"github.com/hashicorp/go-getter"
	logger "github.com/sirupsen/logrus"
)

// Update downloads a github release if its not present in the local config folder
func (release *opsCliRelease) Update() error {

	url, update := release.getLatestDownloadUrl()
	if !update {
		return nil
	}

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
