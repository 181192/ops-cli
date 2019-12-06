package util

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// UsageErrorf shows usage message of cli
func UsageErrorf(cmd *cobra.Command, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s\nSee '%s -h' for help and examples", msg, cmd.CommandPath())
}

// ExecuteCmd execute a command and print result to stdout or stderr
func ExecuteCmd(cmd *cobra.Command, command string, args []string) {
	result := exec.Command(command, args...)
	out, err := result.Output()
	if err != nil {
		fmt.Fprintln(cmd.OutOrStderr(), err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), string(out))
}

// RequireFile checks for a given file and returns an error message to user if it doesn't exists
func RequireFile(fileName string) error {
	if err := checkFile(fileName); err != nil {
		return fmt.Errorf("%s\nRun 'ops download' to download external binaries", err)
	}
	return nil
}

func checkFile(fileName string) error {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return err
	}
	return nil
}
