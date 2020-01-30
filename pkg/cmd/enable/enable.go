package enable

import (
	"time"

	"github.com/181192/ops-cli/pkg/cmd/cmdutils"

	"github.com/spf13/cobra"
)

var (
	kubeconfig        string
	configContext     string
	clusterConfigFile string
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
	profileRevision = "revision"
)

// enableCmd represents the enable command
var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable resource(s)",
	Long:  `Enable resource(s).`,
}

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

// Command will create the `enable` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {

	cmdutils.AddResourceCmd(flagGrouping, enableCmd, enableRepoCmd)

	enableCmd.AddCommand(enableProfileCmd())
	return enableCmd
}

func configureGitOptions(cmd *cobra.Command, opts *GitOptions) *cobra.Command {
	cmd.Flags().StringVar(&opts.URL, gitURL, "", "SSH URL of the Git repository to be used for GitOps, e.g. git@github.com:<github_org>/<repo_name>")
	cmd.Flags().StringVar(&opts.Branch, gitBranch, "", "Git branch to be used for GitOps")
	cmd.Flags().StringVar(&opts.User, gitUser, "", "Username to use as Git committer")
	cmd.Flags().StringVar(&opts.Email, gitEmail, "", "Email to use as Git committer")
	cmd.Flags().StringVar(&opts.PrivateSSHKeyPath, gitPrivateSSHKeyPath, "", "Optional path to the private SSH key to use with Git, e.g. ~/.ssh/id_rsa")

	return cmd
}

func configureFluxOptions(cmd *cobra.Command, opts *FluxOptions) *cobra.Command {
	cmd.Flags().StringSliceVar(&opts.GitPaths, gitPaths, []string{}, "Relative paths within the Git repo for Flux to locate Kubernetes manifests")
	cmd.Flags().StringVar(&opts.GitLabel, gitLabel, "", "Git label to keep track of Flux's sync progress; this is equivalent to overriding --git-sync-tag and --git-notes-ref in Flux")
	cmd.Flags().StringVar(&opts.GitFluxPath, gitFluxPath, "flux/", "Directory within the Git repository where to commit the Flux manifests")
	cmd.Flags().StringVar(&opts.Namespace, namespace, "", "Cluster namespace where to install Flux, the Helm Operator and Tiller")
	cmd.Flags().BoolVar(&opts.WithHelm, withHelm, true, "Install the Helm Operator and Tiller")

	return cmd
}
