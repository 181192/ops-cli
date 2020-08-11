package update

import (
	"net/http"

	"github.com/181192/ops-cli/pkg/util/stringutils"
	"github.com/181192/ops-cli/pkg/util/version"
	semver "github.com/hashicorp/go-version"
	logger "github.com/sirupsen/logrus"
)

func (release *opsCliRelease) getLatestDownloadURL() (string, bool) {

	if version.Version == "unversioned" {
		logger.Warnf("Using non released version, setting version to 0.0 for testing purposes")
		version.Version = "0.0"
	}

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
		return "", false
	}

	url := "https://github.com/" + release.Account + "/" + release.Name + "/releases/download/" + release.Version + "/" + release.ArtifactName

	return url, true
}
