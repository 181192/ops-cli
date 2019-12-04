package util

import (
	"fmt"
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
