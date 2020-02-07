package gitfindroot

import (
	"os/exec"
	"strings"
)

// Stat represent name of git repo and full path
type Stat struct {
	Name string
	Path string
}

// Repo uses git via the console to locate the top level directory
func Repo() (Stat, error) {
	path, err := rootPath()
	if err != nil {
		return Stat{
			Name: "Unknown",
			Path: "./",
		}, err
	}

	gitRepo, err := exec.Command("basename", path).Output()
	if err != nil {
		return Stat{}, err
	}

	return Stat{
		Name: strings.TrimSpace(string(gitRepo)),
		Path: path,
	}, nil
}

func rootPath() (string, error) {
	path, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(path)), nil
}
