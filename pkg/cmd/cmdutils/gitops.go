package cmdutils

import (
	"fmt"
	"time"

	"github.com/181192/ops-cli/pkg/git"
	"github.com/181192/ops-cli/pkg/gitops/profile"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	gitURL               = "git-url"
	gitBranch            = "git-branch"
	gitUser              = "git-user"
	gitEmail             = "git-email"
	gitPrivateSSHKeyPath = "git-private-ssh-key-path"

	gitPaths    = "git-paths"
	gitFluxPath = "git-flux-subdir"
	gitLabel    = "git-label"
	namespace   = "namespace"
	withHelm    = "with-helm"

	profileName     = "name"
	profileOverlay  = "overlay"
	profileRevision = "revision"
)

// InstallOpts are the installation options for Flux
type InstallOpts struct {
	GitOptions  git.Options
	GitPaths    []string
	GitLabel    string
	GitFluxPath string
	Namespace   string
	Timeout     time.Duration
	Amend       bool
	WithHelm    bool
}

// AddCommonFlagsForFlux configures the flags required to install Flux on a
// cluster and have it point to the specified Git repository.
func AddCommonFlagsForFlux(fs *pflag.FlagSet, opts *InstallOpts) {
	AddCommonFlagsForGit(fs, &opts.GitOptions)

	fs.StringSliceVar(&opts.GitPaths, gitPaths, []string{},
		"Relative paths within the Git repo for Flux to locate Kubernetes manifests")
	fs.StringVar(&opts.GitLabel, gitLabel, "flux",
		"Git label to keep track of Flux's sync progress; this is equivalent to overriding --git-sync-tag and --git-notes-ref in Flux")
	fs.StringVar(&opts.GitFluxPath, gitFluxPath, "manifests-flux/",
		"Directory within the Git repository where to commit the Flux manifests")
	fs.StringVar(&opts.Namespace, namespace, "flux-system",
		"Cluster namespace where to install Flux, the Helm Operator and Tiller")
	fs.BoolVar(&opts.WithHelm, withHelm, true,
		"Install the Helm Operator and Tiller")
	fs.BoolVar(&opts.Amend, "amend", false,
		"Stop to manually tweak the Flux manifests before pushing them to the Git repository")
}

// AddCommonFlagsForGit configures the flags required to interact with a Git
// repository.
func AddCommonFlagsForGit(fs *pflag.FlagSet, opts *git.Options) {
	fs.StringVar(&opts.URL, gitURL, "",
		"SSH URL of the Git repository to be used for GitOps, e.g. git@github.com:<github_org>/<repo_name>")
	fs.StringVar(&opts.Branch, gitBranch, "master",
		"Git branch to be used for GitOps")
	fs.StringVar(&opts.User, gitUser, "",
		"Username to use as Git committer")
	fs.StringVar(&opts.Email, gitEmail, "",
		"Email to use as Git committer")
	fs.StringVar(&opts.PrivateSSHKeyPath, gitPrivateSSHKeyPath, "",
		"Optional path to the private SSH key to use with Git, e.g. ~/.ssh/id_rsa")
}

// ValidateGitOptions validates the provided Git options.
func ValidateGitOptions(opts *git.Options) error {
	if err := opts.ValidateURL(); err != nil {
		return errors.Wrapf(err, "please supply a valid --%s argument", gitURL)
	}
	if err := opts.ValidateEmail(); err != nil {
		return fmt.Errorf("please supply a valid --%s argument", gitEmail)
	}
	if err := opts.ValidatePrivateSSHKeyPath(); err != nil {
		return errors.Wrapf(err, "please supply a valid --%s argument", gitPrivateSSHKeyPath)
	}
	return nil
}

// AddCommonFlagsForProfile configures the flags required to enable a Quick
// Start profile.
func AddCommonFlagsForProfile(fs *pflag.FlagSet, opts *profile.Options) {
	fs.StringVarP(&opts.Name, profileName, "", "", "Name or URL of the Quick Start profile. For example, app-dev.")
	fs.StringVarP(&opts.Overlay, profileOverlay, "", "nginx", "Name of the overlay profile. For example nginx,linkerd or istio.")
	fs.StringVarP(&opts.Revision, profileRevision, "", "master", "Revision of the Quick Start profile.")
}

// gitOpsConfigLoader handles loading of ClusterConfigFile v.s. using CLI
// flags for GitOps-related commands.
type gitOpsConfigLoader struct {
	cmd                                *Cmd
	flagsIncompatibleWithConfigFile    sets.String
	flagsIncompatibleWithoutConfigFile sets.String
	validateWithConfigFile             func() error
	validateWithoutConfigFile          func() error
}

// NewGitOpsConfigLoader creates a new ClusterConfigLoader which handles
// loading of ClusterConfigFile v.s. using CLI flags for GitOps-related
// commands.
func NewGitOpsConfigLoader(cmd *Cmd) ClusterConfigLoader {
	l := &gitOpsConfigLoader{
		cmd: cmd,
		flagsIncompatibleWithConfigFile: sets.NewString(
			"location",
			"version",
			"cluster",
		),
		flagsIncompatibleWithoutConfigFile: sets.NewString(),
	}

	l.validateWithoutConfigFile = func() error {
		config := l.cmd.ClusterConfig
		if config.ObjectMeta.Name == "" {
			return ErrMustBeSet(ClusterNameFlag(cmd))
		}
		if config.Spec.Location == "" {
			return ErrMustBeSet("--location")
		}
		return nil
	}

	l.validateWithConfigFile = func() error {
		config := l.cmd.ClusterConfig
		if config.ObjectMeta.Name == "" {
			return ErrMustBeSet("metadata.name")
		}

		if config.Spec.Location == "" {
			return ErrMustBeSet("spec.location")
		}
		return nil
	}

	return l
}

// Load ClusterConfig or use CLI flags.
func (l *gitOpsConfigLoader) Load() error {
	if l.cmd.ClusterConfigFile == "" {
		for f := range l.flagsIncompatibleWithoutConfigFile {
			if flag := l.cmd.CobraCommand.Flag(f); flag != nil && flag.Changed {
				return fmt.Errorf("cannot use --%s unless a config file is specified via --config-file/-f", f)
			}
		}
		return l.validateWithoutConfigFile()
	}

	var err error

	// The reference to ClusterConfig should only be reassigned if ClusterConfigFile is specified
	// because other parts of the code store the pointer locally and access it directly instead of via
	// the Cmd reference
	if l.cmd.ClusterConfig, err = LoadConfigFromFile(l.cmd.ClusterConfigFile); err != nil {
		return err
	}

	for f := range l.flagsIncompatibleWithConfigFile {
		if flag := l.cmd.CobraCommand.Flag(f); flag != nil && flag.Changed {
			return ErrCannotUseWithConfigFile(fmt.Sprintf("--%s", f))
		}
	}

	return l.validateWithConfigFile()
}
