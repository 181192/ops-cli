package gitconfig

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"syscall"
)

// ErrNotFound error not found
type ErrNotFound struct {
	Key string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("the key `%s` is not found", e.Key)
}

// Entire extracts configuration value from `$HOME/.gitconfig` file ,
// `$GIT_CONFIG`, /etc/gitconfig or include.path files.
func Entire(key string) (string, error) {
	return execGitConfig(key)
}

// Username extracts git user name from `Entire gitconfig`.
// This is same as Entire("user.name")
func Username() (string, error) {
	return Entire("user.name")
}

// Email extracts git user email from `$HOME/.gitconfig` file or `$GIT_CONFIG`.
// This is same as Global("user.email")
func Email() (string, error) {
	return Entire("user.email")
}

func execGitConfig(args ...string) (string, error) {
	gitArgs := append([]string{"config", "--get", "--null"}, args...)
	var stdout bytes.Buffer
	cmd := exec.Command("git", gitArgs...)
	cmd.Stdout = &stdout
	cmd.Stderr = ioutil.Discard

	err := cmd.Run()
	if exitError, ok := err.(*exec.ExitError); ok {
		if waitStatus, ok := exitError.Sys().(syscall.WaitStatus); ok {
			if waitStatus.ExitStatus() == 1 {
				return "", &ErrNotFound{Key: args[len(args)-1]}
			}
		}
		return "", err
	}

	return strings.TrimRight(stdout.String(), "\000"), nil
}
