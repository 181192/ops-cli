package cmdutils

import (
	"fmt"
	"os"
	"strings"

	api "github.com/181192/ops-cli/pkg/apis/opscli.io/v1alpha1"

	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

// IncompatibleFlags is a common substring of an error message
const IncompatibleFlags = "cannot be used at the same time"

// GetNameArg tests to ensure there is only 1 name argument
func GetNameArg(args []string) string {
	if len(args) > 1 {
		logger.Fatal("only one argument is allowed to be used as a name")
		os.Exit(1)
	}
	if len(args) == 1 {
		return (strings.TrimSpace(args[0]))
	}
	return ""
}

// AddLocationFlag adds common --location flag
func AddLocationFlag(fs *pflag.FlagSet, p *api.AKSClusterConfigSpec) {
	fs.StringVarP(&p.Location, "location", "l", "", "Azure Location")
}

// AddVersionFlag adds common --version flag
func AddVersionFlag(fs *pflag.FlagSet, meta *api.AKSClusterConfig, extraUsageInfo string) {
	usage := fmt.Sprintf("Kubernetes version (valid options: %s)", strings.Join(api.SupportedVersions(), ", "))
	if extraUsageInfo != "" {
		usage = fmt.Sprintf("%s [%s]", usage, extraUsageInfo)
	}
	fs.StringVar(&meta.Spec.Version, "version", meta.Spec.Version, usage)
}

// AddWaitFlag adds common --wait flag
func AddWaitFlag(fs *pflag.FlagSet, wait *bool, description string) {
	AddWaitFlagWithFullDescription(fs, wait, fmt.Sprintf("wait for %s before exiting", description))
}

// AddWaitFlagWithFullDescription adds common --wait flag
func AddWaitFlagWithFullDescription(fs *pflag.FlagSet, wait *bool, description string) {
	fs.BoolVarP(wait, "wait", "w", *wait, description)
}

// AddCommonFlagsForKubeconfig adds common flags for controlling how output kubeconfig is written
func AddCommonFlagsForKubeconfig(fs *pflag.FlagSet, outputPath, authenticatorRoleARN *string, setContext, autoPath *bool, exampleName string) {
	fs.StringVar(outputPath, "kubeconfig", "", "path to write kubeconfig")
	fs.BoolVar(setContext, "set-kubeconfig-context", true, "if true then current-context will be set in kubeconfig; if a context is already set then it will be overwritten")
}

// ErrClusterFlagAndArg wraps ErrFlagAndArg() by passing in the
// proper flag name.
func ErrClusterFlagAndArg(cmd *Cmd, nameFlag, nameArg string) error {
	return ErrFlagAndArg(ClusterNameFlag(cmd), nameFlag, nameArg)
}

// ErrFlagAndArg may be used to err for options that can be given
// as flags /and/ arg but only one is allowed to be used.
func ErrFlagAndArg(kind, flag, arg string) error {
	return fmt.Errorf("%s=%s and argument %s %s", kind, flag, arg, IncompatibleFlags)
}

// ErrMustBeSet is a common error message
func ErrMustBeSet(pathOrFlag string) error {
	return fmt.Errorf("%s must be set", pathOrFlag)
}

// ErrCannotUseWithConfigFile is a common error message
func ErrCannotUseWithConfigFile(what string) error {
	return fmt.Errorf("cannot use %s when --config-file/-f is set", what)
}

// ErrUnsupportedManagedFlag reports unsupported flags for Managed Nodegroups
func ErrUnsupportedManagedFlag(flag string) error {
	return fmt.Errorf("%s is not supported for Managed Nodegroups (--managed=true)", flag)
}

// ErrUnsupportedNameArg reports unsupported usage of `name` argument
func ErrUnsupportedNameArg() error {
	return errors.New("name argument is not supported")
}

// ClusterNameFlag returns the flag to use for the cluster name
// taking the principal resource into account.
func ClusterNameFlag(cmd *Cmd) string {
	if cmd.CobraCommand.Use == "cluster" {
		return "--name"
	}
	return "--cluster"
}
