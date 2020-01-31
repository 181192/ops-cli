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

	ClusterConfig *api.ClusterConfig
}

// NewCtl validates
func (c *Cmd) NewCtl() error {
	if err := api.ValidateClusterConfig(c.ClusterConfig); err != nil {
		if c.Validate {
			return err
		}
		logger.Warningf("ignoring validation error: %s", err.Error())
	}

	return nil
}

// AddResourceCmd create a registers a new command under the given verb command
func AddResourceCmd(flagGrouping *FlagGrouping, parentVerbCmd *cobra.Command, newCmd func(*Cmd)) {
	c := &Cmd{
		CobraCommand: &cobra.Command{},
		Validate:     true,
	}
	c.FlagSetGroup = flagGrouping.New(c.CobraCommand)
	newCmd(c)
	c.FlagSetGroup.AddTo(c.CobraCommand)
	parentVerbCmd.AddCommand(c.CobraCommand)
}
