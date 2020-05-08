package download

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/181192/ops-cli/pkg/download"
	"github.com/181192/ops-cli/pkg/util"
	"github.com/hashicorp/go-getter"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

var terraformBinary = util.GetConfigDirectory() + "/bin/terraform"

type terraformRelease struct {
	*download.Release
}

// https://releases.hashicorp.com/terraform/0.12.17/terraform_0.12.17_linux_amd64.zip
// https://releases.hashicorp.com/terraform/0.12.17/terraform_0.12.17_darwin_amd64.zip
// https://releases.hashicorp.com/terraform/0.12.17/terraform_0.12.17_windows_amd64.zip

func newTerraformRelease() *terraformRelease {

	var artifactsNames = map[string]string{
		"darwin/amd64":  "darwin_amd64.zip",
		"linux/amd64":   "linux_amd64.zip",
		"windows/amd64": "windows_amd64.zip",
	}

	release := &terraformRelease{
		&download.Release{
			Name:          "terraform",
			LocalFileName: terraformBinary,
		},
	}

	release.SetVersion()
	release.SetArtifactName(artifactsNames)
	release.setDownloadURL()

	return release
}

// terraformCmd represents the terraform command
var terraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "Terraform IaC tool",
	Long:  `Terraform IaC tool.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return util.RequireFile(terraformBinary)
	},
	Run: func(cmd *cobra.Command, args []string) {
		util.ExecuteCmd(cmd, terraformBinary, args)
	},
}

func (release *terraformRelease) setDownloadURL() *terraformRelease {
	if release.Version == "latest" {
		terraformVersion := "0.12.17"

		resp, err := http.Get("https://checkpoint-api.hashicorp.com/v1/check/terraform")
		if err != nil {
			logger.Warnf("Failed to get latest stable version of terraform %s %s", err, resp.Status)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			logger.Warnf("Failed to read latest version, default to %s", terraformVersion)
			release.Version = terraformVersion
		}

		latestTag := gjson.Get(string(body), "current_version").String()
		logger.Debugf("Latest relese version of terraform %s", latestTag)
		release.Version = latestTag
	}
	release.URL = "https://releases.hashicorp.com/terraform/" + release.Version + "/terraform_" + release.Version + "_" + release.ArtifactName
	return release
}

func (release *terraformRelease) DownloadIfNotExists() error {
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
