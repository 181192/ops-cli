package file

import (
	"fmt"
	"io"
	"os"

	kopsutils "k8s.io/kops/upup/pkg/fi/utils"
)

// Exists checks to see if a file exists.
func Exists(path string) bool {
	extendedPath := ExpandPath(path)
	_, err := os.Stat(extendedPath)
	return err == nil
}

// ExpandPath expands path with ~ notation
func ExpandPath(p string) string { return kopsutils.ExpandPath(p) }

// Copy copies a file from src to dst
func Copy(srcFile, dstFile string) error {
	out, err := os.Create(dstFile)
	defer out.Close()
	if err != nil {
		return err
	}

	in, err := os.Open(srcFile)
	defer in.Close()
	if err != nil {
		return err
	}

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

// CreateIfNotExists creates a file if not exists
func CreateIfNotExists(dir string, perm os.FileMode) error {
	if Exists(dir) {
		return nil
	}

	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}

// CopySymLink copies symlink from src to dst
func CopySymLink(source, dest string) error {
	link, err := os.Readlink(source)
	if err != nil {
		return err
	}
	return os.Symlink(link, dest)
}
