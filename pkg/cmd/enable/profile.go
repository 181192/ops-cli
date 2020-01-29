package enable

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"

	"github.com/181192/ops-cli/pkg/gitops"
	"github.com/181192/ops-cli/pkg/gitops/fileprocessor"

	api "github.com/181192/ops-cli/pkg/apis/opscli.io/v1alpha1"
	"github.com/181192/ops-cli/pkg/git"
	"github.com/181192/ops-cli/pkg/gitops/profile"

	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
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

func enableProfileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Set up Flux and deploy the components from the selected profile",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			// client, err := kubernetes.NewClient(kubeconfig, configContext)
			// if err != nil {
			// 	return fmt.Errorf("failed to create k8s client: %v", err)
			// }

			// opts := configureProfileCmd(cmd)

			profileName := getNameArg(args)

			opts := &ProfileOptions{
				profileOptions: profile.Options{
					Name: profileName,
				},
				gitOptions: git.Options{
					URL: profileName,
				},
			}

			return doEnableProfile(cmd, opts)
		},
	}

	return cmd
}

func doEnableProfile(cmd *cobra.Command, opts *ProfileOptions) error {
	if err := opts.Validate(); err != nil {
		return err
	}

	profileRepoURL, err := profile.RepositoryURL(opts.profileOptions.Name)
	if err != nil {
		return errors.Wrap(err, "please supply a valid profile name or URL")
	}

	// Load GitOpsConfig
	// if err := NewGitOpsConfigLoader(cmd).Load(); err != nil {
	// 	return err
	// }

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

	clusterConfig := api.NewAKSClusterConfig("", "", api.AKSClusterConfig{})
	clusterConfig.Name = "test"

	profile := &gitops.Profile{
		Processor: &fileprocessor.GoTemplateProcessor{
			Params: fileprocessor.NewTemplateParameters(clusterConfig),
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
