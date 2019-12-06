package github

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
}

// getDownloadURL based on if version is latest or tagged
func getDownloadURL(release Release) string {
	if release.Version == "latest" {
		return "https://github.com/" + release.Account + "/" + release.Name + "/releases/latest/download/" + release.ArtifactName
	}
	return "https://github.com/" + release.Account + "/" + release.Name + "/releases/download/" + release.Version + "/" + release.ArtifactName
}

// DownloadIfNotExists downloads a github release if its not present in the local config folder
func DownloadIfNotExists(release Release) error {
	if _, err := os.Stat(release.LocalFileName); os.IsNotExist(err) {
		progress := getter.WithProgress(DefaultProgressBar)
		url := getDownloadURL(release)

		fmt.Println("Attempting to download helmfile, version: " + release.Version)
		err := getter.GetFile(release.LocalFileName, url, progress)
		if err != nil {
			return fmt.Errorf("%s\nFailed to to download external binaries", err)
		}
		os.Chmod(release.LocalFileName, 0775)
	}
	return nil
}

func GetVersion(release Release) string {
	version := viper.GetString("apps. " + release.Name)
	if version == "" {
		fmt.Println("No version set in config, using latest")
		return "latest"
	}
	return version
}

func GetArtifactName(artifacts map[string]string) string {
	platform := runtime.GOOS + "/" + runtime.GOARCH
	artifact, found := artifacts[platform]
	if !found {
		fmt.Printf("Unsupported os/platform to use helmfile %s\n", platform)
	}
	return artifact
}
