package generate

import (
	"context"
	"os"
	"strings"

	api "github.com/181192/ops-cli/pkg/apis/opscli.io/v1alpha1"
	"github.com/181192/ops-cli/pkg/cmd/cmdutils"
	"github.com/181192/ops-cli/pkg/generated/clientset/versioned/scheme"
	"github.com/181192/ops-cli/pkg/git"
	"github.com/181192/ops-cli/pkg/gitops"
	"github.com/181192/ops-cli/pkg/gitops/fileprocessor"
	"github.com/181192/ops-cli/pkg/gitops/profile"

	"github.com/kr/pretty"
	"github.com/pkg/errors"
	"github.com/rancher/wrangler/pkg/schemes"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// ProfileOptions groups input for the "enable profile" command
type ProfileOptions struct {
	gitOptions     git.Options
	profileOptions profile.Options
}

// Validate validates this ProfileOptions object
func (opts ProfileOptions) Validate() error {
	if err := opts.gitOptions.Validate(); err != nil {
		return err
	}
	return opts.profileOptions.Validate()
}

// GetNameArg tests to ensure there is only 1 name argument
func getNameArg(args []string) string {
	if len(args) > 1 {
		logger.Fatal("only one argument is allowed to be used as a name")
		os.Exit(1)
	}
	if len(args) == 1 {
		return (strings.TrimSpace(args[0]))
	}
	return ""
}

// Options options for generating profile
type Options struct {
	GitOptions        git.Options
	ProfilePath       string
	PrivateSSHKeyPath string
}

const (
	gitURL      = "git-url"
	gitBranch   = "git-branch"
	profilePath = "profile-path"
)

func generateProfileCmd() *cobra.Command {

	var opts Options

	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Generate a gitops profile",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			scheme.AddToScheme(schemes.All)

			// profileName := getNameArg(args)

			return doGenerateProfile(cmd, opts)
		},
	}

	cmd.Flags().StringVar(&opts.GitOptions.URL, gitURL, "", "URL for the base repository")
	cmd.Flags().StringVar(&opts.GitOptions.Branch, gitBranch, "master", "Git branch")
	cmd.Flags().StringVar(&opts.ProfilePath, profilePath, "./", "Path to generate the profile in")
	cobra.MarkFlagRequired(cmd.Flags(), gitURL)

	return cmd
}

func doGenerateProfile(cmd *cobra.Command, opts Options) error {
	// if err := opts.Validate(); err != nil {
	// 	return err
	// }
	api.NewAKSClusterConfig("", "", api.AKSClusterConfig{})

	clusterConfig, err := cmdutils.LoadConfigFromFile(clusterConfigFile)
	if err != nil {
		return err
	}

	logger.Debugf("%# v", pretty.Formatter(clusterConfig))

	processor := &fileprocessor.GoTemplateProcessor{
		Params: fileprocessor.NewTemplateParameters(clusterConfig),
	}

	profile := &gitops.Profile{
		Processor: processor,
		Path:      opts.ProfilePath,
		GitOpts:   opts.GitOptions,
		GitCloner: git.NewGitClient(git.ClientParams{
			PrivateSSHKeyPath: opts.PrivateSSHKeyPath,
		}),
		FS: afero.NewOsFs(),
		IO: afero.Afero{Fs: afero.NewOsFs()},
	}

	err = profile.Generate(context.Background())
	if err != nil {
		return errors.Wrap(err, "error generating profile")
	}

	profile.DeleteClonedDirectory()
	return nil
}
