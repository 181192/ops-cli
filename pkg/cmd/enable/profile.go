package enable

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	api "github.com/181192/ops-cli/pkg/apis/opscli.io/v1alpha1"
	"github.com/181192/ops-cli/pkg/cmd/cmdutils"
	"github.com/181192/ops-cli/pkg/git"
	"github.com/181192/ops-cli/pkg/git/gitconfig"
	"github.com/181192/ops-cli/pkg/gitops"
	"github.com/181192/ops-cli/pkg/gitops/fileprocessor"
	"github.com/181192/ops-cli/pkg/gitops/profile"

	"github.com/kr/pretty"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// ProfileOptions groups input for the "enable profile" command
type ProfileOptions struct {
	gitOptions     git.Options
	profileOptions profile.Options
}

func enableProfileCmd(cmd *cmdutils.Cmd) {
	var opts ProfileOptions

	cmd.ClusterConfig = api.DefaultAKSClusterConfig()
	cmd.CobraCommand.Use = "profile"
	cmd.CobraCommand.Short = "Enable and deploy the components from the selected profile"
	cmd.CobraCommand.Long = ""
	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doEnableProfile(cmd, &opts)
	}

	cmd.FlagSetGroup.InFlagSet("Enable profile", func(fs *pflag.FlagSet) {
		cmdutils.AddCommonFlagsForProfile(fs, &opts.profileOptions)
		cmdutils.AddCommonFlagsForGit(fs, &opts.gitOptions)
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddLocationFlag(fs, cmd.ClusterConfig)
		cmdutils.AddClusterFlag(fs, cmd.ClusterConfig)
	})
}

func doEnableProfile(cmd *cmdutils.Cmd, opts *ProfileOptions) error {
	if cmd.NameArg != "" && opts.profileOptions.Name != "" {
		return cmdutils.ErrClusterFlagAndArg(cmd, cmd.NameArg, opts.profileOptions.Name)
	}

	if cmd.NameArg != "" {
		opts.profileOptions.Name = cmd.NameArg
	}

	if opts.gitOptions.User == "" {
		if gitUser, err := gitconfig.Username(); err == nil {
			opts.gitOptions.User = gitUser
		}
	}

	if opts.gitOptions.Email == "" {
		if gitEmail, err := gitconfig.Email(); err == nil {
			opts.gitOptions.Email = gitEmail
		}
	}

	logger.Debugf("%# v", pretty.Formatter(opts))

	if err := opts.gitOptions.Validate(); err != nil {
		return err
	}

	if err := opts.profileOptions.Validate(); err != nil {
		return err
	}

	profileRepoURL, err := profile.RepositoryURL(opts.profileOptions.Name)
	if err != nil {
		return errors.Wrap(err, "please supply a valid profile name or URL")
	}

	if err := cmdutils.NewGitOpsConfigLoader(cmd).Load(); err != nil {
		return err
	}

	// Clone user's repo to apply profile
	usersRepoName, err := git.RepoName(opts.gitOptions.URL)
	if err != nil {
		return err
	}

	usersRepoDir, err := ioutil.TempDir("", usersRepoName+"-")
	logger.Debugf("Directory %s will be used to clone the configuration repository and install the profile", usersRepoDir)
	profileOutputPath := filepath.Join(usersRepoDir, "base")

	gitClient := git.NewGitClient(git.ClientParams{
		PrivateSSHKeyPath: opts.gitOptions.PrivateSSHKeyPath,
	})

	err = gitClient.CloneRepoInPath(
		usersRepoDir,
		git.CloneOptions{
			URL:       opts.gitOptions.URL,
			Branch:    opts.gitOptions.Branch,
			Bootstrap: true,
		},
	)
	if err != nil {
		return err
	}

	profile := &gitops.Profile{
		Processor: &fileprocessor.GoTemplateProcessor{
			Params: fileprocessor.NewTemplateParameters(cmd.ClusterConfig),
		},
		Path: profileOutputPath,
		GitOpts: git.Options{
			URL:    profileRepoURL,
			Branch: opts.profileOptions.Revision,
		},
		GitCloner: git.NewGitClient(git.ClientParams{
			PrivateSSHKeyPath: opts.gitOptions.PrivateSSHKeyPath,
		}),
		FS: afero.NewOsFs(),
		IO: afero.Afero{Fs: afero.NewOsFs()},
	}

	err = profile.Generate(context.Background())
	if err != nil {
		return errors.Wrap(err, "error generating profile")
	}

	// Git add, commit and push component files in the user's repo
	if err = gitClient.Add("."); err != nil {
		return err
	}

	commitMsg := fmt.Sprintf("Add %s profile components", opts.profileOptions.Name)
	if err = gitClient.Commit(commitMsg, opts.gitOptions.User, opts.gitOptions.Email); err != nil {
		return err
	}

	// if err = gitClient.Push(); err != nil {
	// 	return err
	// }

	profile.DeleteClonedDirectory()
	os.RemoveAll(usersRepoDir)

	return nil
}
