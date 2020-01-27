package cmdutils

import (
	api "github.com/181192/ops-cli/pkg/apis/opscli.io/v1alpha1"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Cmd hold common attributes between commands
type Cmd struct {
	CobraCommand *cobra.Command
	FlagSetGroup *NamedFlagSetGroup

	Validate bool

	NameArg           string
	ClusterConfigFile string

	ClusterConfig *api.AKSClusterConfig
}

// NewCtl validates
func (c *Cmd) NewCtl() error {
	if err := api.ValidateAKSClusterConfig(c.ClusterConfig); err != nil {
		if c.Validate {
			return err
		}
		logger.Warningf("ignoring validation error: %s", err.Error())
	}

	return nil
}
