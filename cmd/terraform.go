package cmd

import (
	"fmt"
	"os"

	"github.com/181192/ops-cli/pkg/download"
	cmdUtil "github.com/181192/ops-cli/pkg/util"
	"github.com/hashicorp/go-getter"
	"github.com/spf13/cobra"
)

var terraformBinary = cfgFolder + "/bin/terraform"

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
		return cmdUtil.RequireFile(terraformBinary)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmdUtil.ExecuteCmd(cmd, terraformBinary, args)
	},
}

// https://releases.hashicorp.com/terraform/{{ version }}/terraform_{{ version }}_linux_amd64.zip
// https://releases.hashicorp.com/terraform/0.12.17/terraform_0.12.17_linux_amd64.zip
// https://releases.hashicorp.com/terraform/0.12.17/terraform_0.12.17_darwin_amd64.zip
// https://releases.hashicorp.com/terraform/0.12.17/terraform_0.12.17_windows_amd64.zip

func init() {
	rootCmd.AddCommand(terraformCmd)
}

func (release *terraformRelease) setDownloadURL() *terraformRelease {
	if release.Version == "latest" {
		release.Version = "0.12.17"
	}
	release.URL = "https://releases.hashicorp.com/terraform/" + release.Version + "/terraform_" + release.Version + "_" + release.ArtifactName
	return release
}

func (release *terraformRelease) DownloadIfNotExists() error {
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
