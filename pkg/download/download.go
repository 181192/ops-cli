package download

import (
	"fmt"
	"os"
	"runtime"

	"github.com/hashicorp/go-getter"
	"github.com/spf13/viper"
)

// Release represents a github release
type Release struct {
	Account       string
	Name          string
	ArtifactName  string
	LocalFileName string
	Version       string
	URL           string
}

// SetDownloadURL based on if version is latest or tagged
func (release *Release) SetDownloadURL() *Release {
	url := "https://github.com/" + release.Account + "/" + release.Name + "/releases/download/" + release.Version + "/" + release.ArtifactName

	if release.Version == "latest" {
		url = "https://github.com/" + release.Account + "/" + release.Name + "/releases/latest/download/" + release.ArtifactName
	}

	release.URL = url
	return release
}

// DownloadIfNotExists downloads a github release if its not present in the local config folder
func (release *Release) DownloadIfNotExists() error {
	if _, err := os.Stat(release.LocalFileName); os.IsNotExist(err) {
		progress := getter.WithProgress(DefaultProgressBar)

		fmt.Printf("Attempting to download %s, version %s, to %q\n", release.Name, release.Version, release.LocalFileName)
		err := getter.GetFile(release.LocalFileName, release.URL, progress)
		if err != nil {
			return fmt.Errorf("%s\nFailed to to download external binaries", err)
		}
		os.Chmod(release.LocalFileName, 0775)
	} else {
		fmt.Printf("%s already exists at %s\n", release.Name, release.LocalFileName)
	}
	return nil
}

// SetVersion return version from config file or latest if not specified
func (release *Release) SetVersion() *Release {
	version := viper.GetString("apps. " + release.Name)
	if version == "" {
		fmt.Println("No version set in config, using latest")
		version = "latest"
	}

	release.Version = version
	return release
}

// SetArtifactName returns artifactname given GOOS and GOARCH
func (release *Release) SetArtifactName(artifacts map[string]string) *Release {
	platform := runtime.GOOS + "/" + runtime.GOARCH
	artifact, found := artifacts[platform]
	if !found {
		fmt.Printf("Unsupported os/platform %s\n", platform)
	}

	release.ArtifactName = artifact
	return release
}
