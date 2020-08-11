// +build windows

package update

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"

	"github.com/181192/ops-cli/pkg/download"
	"github.com/hashicorp/go-getter"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/sys/windows"
)

// Update downloads a github release if its not present in the local config folder
func (release *opsCliRelease) Update() error {

	url, update := release.getLatestDownloadURL()
	if !update {
		return nil
	}

	if !isWinAdmin() {
		if err := runElevated(); err != nil {
			return err
		}
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

func isWinAdmin() bool {
	if _, err := os.Open("\\\\.\\PHYSICALDRIVE0"); err != nil {
		return false
	}
	return true
}

func runElevated() error {
	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	args := strings.Join(os.Args[1:], " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)

	var showCmd int32 = 1 //SW_NORMAL

	err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		return fmt.Errorf("%s\nFailed to exec as elevated user", err)
	}

	return nil
}
