package generate

import (
	"context"

	api "github.com/181192/ops-cli/pkg/apis/opscli.io/v1alpha1"
	"github.com/181192/ops-cli/pkg/cmd/cmdutils"
	"github.com/181192/ops-cli/pkg/generated/clientset/versioned/scheme"
	"github.com/181192/ops-cli/pkg/git"
	"github.com/181192/ops-cli/pkg/gitops"
	"github.com/181192/ops-cli/pkg/gitops/fileprocessor"

	"github.com/kr/pretty"
	"github.com/pkg/errors"
	"github.com/rancher/wrangler/pkg/schemes"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	gitURL               = "git-url"
	gitBranch            = "git-branch"
	profilePath          = "profile-path"
	gitPrivateSSHKeyPath = "git-private-ssh-key-path"
)

// ProfileOptions options for generating profile
type ProfileOptions struct {
	GitOptions        git.Options
	ProfilePath       string
	PrivateSSHKeyPath string
}

func generateProfileCmd(cmd *cmdutils.Cmd) {
	var opts ProfileOptions

	cmd.ClusterConfig = api.DefaultAKSClusterConfig()
	cmd.CobraCommand.Use = "profile"
	cmd.CobraCommand.Short = "Generate a gitops profile"
	cmd.CobraCommand.Long = ""
	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doGenerateProfile(cmd, opts)
	}

	cmd.FlagSetGroup.InFlagSet("Generate Profile", func(fs *pflag.FlagSet) {
		fs.StringVarP(&opts.GitOptions.URL, gitURL, "", "", "URL for the quickstart base repository")
		fs.StringVarP(&opts.GitOptions.Branch, gitBranch, "", "master", "Git branch")
		fs.StringVar(&opts.PrivateSSHKeyPath, gitPrivateSSHKeyPath, "", "Optional path to the private SSH key to use with Git, e.g. ~/.ssh/id_rsa")
		fs.StringVarP(&opts.ProfilePath, profilePath, "", "./", "Path to generate the profile in")
		_ = cobra.MarkFlagRequired(fs, gitURL)
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddLocationFlag(fs, cmd.ClusterConfig)
		cmdutils.AddClusterFlag(fs, cmd.ClusterConfig)
	})
}

func doGenerateProfile(cmd *cmdutils.Cmd, opts ProfileOptions) error {
	scheme.AddToScheme(schemes.All)

	if err := cmdutils.NewMetadataLoader(cmd).Load(); err != nil {
		return err
	}

	logger.Debugf("%# v", pretty.Formatter(cmd.ClusterConfig))

	processor := &fileprocessor.GoTemplateProcessor{
		Params: fileprocessor.NewTemplateParameters(cmd.ClusterConfig),
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

	if err := profile.Generate(context.Background()); err != nil {
		return errors.Wrap(err, "error generating profile")
	}

	profile.DeleteClonedDirectory()
	return nil
}
