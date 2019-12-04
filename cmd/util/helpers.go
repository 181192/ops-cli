package util

import (
	"fmt"

	"github.com/spf13/cobra"
)

// UsageErrorf shows usage message of cli
func UsageErrorf(cmd *cobra.Command, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s\nSee '%s -h' for help and examples", msg, cmd.CommandPath())
}
