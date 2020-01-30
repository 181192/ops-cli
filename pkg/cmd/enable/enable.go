package enable

import (
	"time"

	"github.com/181192/ops-cli/pkg/cmd/cmdutils"

	"github.com/spf13/cobra"
)

// GitOptions represent git options from cli
type GitOptions struct {
	URL               string
	Branch            string
	User              string
	Email             string
	PrivateSSHKeyPath string
}

// FluxOptions represent flux options from cli
type FluxOptions struct {
	GitOptions  GitOptions
	GitPaths    []string
	GitLabel    string
	GitFluxPath string
	Namespace   string
	Timeout     time.Duration
	Amend       bool // TODO: remove, as we eventually no longer want to support this mode?
	WithHelm    bool
}

// enableCmd represents the enable command
var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable resource(s)",
	Long:  `Enable resource(s).`,
}

// Command will create the `enable` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {

	cmdutils.AddResourceCmd(flagGrouping, enableCmd, enableRepoCmd)
	cmdutils.AddResourceCmd(flagGrouping, enableCmd, enableProfileCmd)

	return enableCmd
}
