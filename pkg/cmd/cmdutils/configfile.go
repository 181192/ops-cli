package cmdutils

import (
	"fmt"

	"github.com/181192/ops-cli/pkg/ops"

	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/sets"
)

// AddConfigFileFlag adds common --config-file flag
func AddConfigFileFlag(fs *pflag.FlagSet, path *string) {
	fs.StringVarP(path, "config-file", "f", "", "load configuration from a file (or stdin if set to '-')")
}

// ClusterConfigLoader is an interface that loaders should implement
type ClusterConfigLoader interface {
	Load() error
}

type commonClusterConfigLoader struct {
	*Cmd

	flagsIncompatibleWithConfigFile, flagsIncompatibleWithoutConfigFile sets.String

	validateWithConfigFile, validateWithoutConfigFile func() error
}

var (
	defaultFlagsIncompatibleWithConfigFile = sets.NewString(
		"name",
		"location",
		"version",
		"cluster",
		"namepace",
	)
	defaultFlagsIncompatibleWithoutConfigFile = sets.NewString(
		"only",
		"include",
		"exclude",
		"only-missing",
	)
)

func newCommonClusterConfigLoader(cmd *Cmd) *commonClusterConfigLoader {
	nilValidatorFunc := func() error { return nil }

	return &commonClusterConfigLoader{
		Cmd: cmd,

		validateWithConfigFile:             nilValidatorFunc,
		flagsIncompatibleWithConfigFile:    defaultFlagsIncompatibleWithConfigFile,
		validateWithoutConfigFile:          nilValidatorFunc,
		flagsIncompatibleWithoutConfigFile: defaultFlagsIncompatibleWithoutConfigFile,
	}
}

// Load ClusterConfig or use flags
func (l *commonClusterConfigLoader) Load() error {

	if l.ClusterConfigFile == "" {
		for f := range l.flagsIncompatibleWithoutConfigFile {
			if flag := l.CobraCommand.Flag(f); flag != nil && flag.Changed {
				return fmt.Errorf("cannot use --%s unless a config file is specified via --config-file/-f", f)
			}
		}
		return l.validateWithoutConfigFile()
	}

	var err error

	// The reference to ClusterConfig should only be reassigned if ClusterConfigFile is specified
	// because other parts of the code store the pointer locally and access it directly instead of via
	// the Cmd reference
	if l.ClusterConfig, err = ops.LoadConfigFromFile(l.ClusterConfigFile); err != nil {
		return err
	}
	meta := l.ClusterConfig.ObjectMeta
	spec := l.ClusterConfig.Spec

	for f := range l.flagsIncompatibleWithConfigFile {
		if flag := l.CobraCommand.Flag(f); flag != nil && flag.Changed {
			return ErrCannotUseWithConfigFile(fmt.Sprintf("--%s", f))
		}
	}

	if l.NameArg != "" {
		return ErrCannotUseWithConfigFile(fmt.Sprintf("name argument %q", l.NameArg))
	}

	if meta.Name == "" {
		return ErrMustBeSet("metadata.name")
	}

	if spec.Location == "" {
		return ErrMustBeSet("spec.location")
	}

	return l.validateWithConfigFile()
}

func (l *commonClusterConfigLoader) validateMetadataWithoutConfigFile() error {
	meta := l.ClusterConfig.ObjectMeta

	if meta.Name != "" && l.NameArg != "" {
		return ErrClusterFlagAndArg(l.Cmd, meta.Name, l.NameArg)
	}

	if l.NameArg != "" {
		meta.Name = l.NameArg
	}

	if meta.Name == "" {
		return ErrMustBeSet(ClusterNameFlag(l.Cmd))
	}

	return nil
}

// NewMetadataLoader handles loading of clusterConfigFile vs using flags for all commands that require only
// metadata fields, e.g. `eksctl delete cluster` or `eksctl utils update-kube-proxy` and other similar
// commands that do simple operations against existing clusters
func NewMetadataLoader(cmd *Cmd) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)

	l.validateWithoutConfigFile = l.validateMetadataWithoutConfigFile

	return l
}
