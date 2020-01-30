package enable

import (
	api "github.com/181192/ops-cli/pkg/apis/opscli.io/v1alpha1"
	"github.com/181192/ops-cli/pkg/cmd/cmdutils"
	scheme "github.com/181192/ops-cli/pkg/generated/clientset/versioned/scheme"
	"github.com/181192/ops-cli/pkg/git"
	"github.com/181192/ops-cli/pkg/helm"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	"github.com/kr/pretty"
	"github.com/rancher/wrangler/pkg/schemes"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func enableRepoCmd(cmd *cmdutils.Cmd) {
	var opts cmdutils.InstallOpts

	cmd.ClusterConfig = api.DefaultAKSClusterConfig()
	cmd.CobraCommand.Use = "repo"
	cmd.CobraCommand.Short = "Set up a repo for gitops, installing Flux in the cluster and initializing its manifests"
	cmd.CobraCommand.Long = ""
	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doEnableRepo(cmd, &opts)
	}

	cmd.FlagSetGroup.InFlagSet("Enable repo", func(fs *pflag.FlagSet) {
		cmdutils.AddCommonFlagsForFlux(fs, &opts)
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddLocationFlag(fs, cmd.ClusterConfig)
		cmdutils.AddClusterFlag(fs, cmd.ClusterConfig)
	})
}

func doEnableRepo(cmd *cmdutils.Cmd, opts *cmdutils.InstallOpts) error {
	scheme.AddToScheme(schemes.All)
	gitOpts := opts.GitOptions

	if err := cmdutils.ValidateGitOptions(&gitOpts); err != nil {
		return err
	}

	if err := cmdutils.NewGitOpsConfigLoader(cmd).Load(); err != nil {
		return err
	}

	logger.Debugf("%# v", pretty.Formatter(cmd.ClusterConfig))
	logger.Debugf("%# v", pretty.Formatter(opts))

	gitClient := git.NewGitClient(git.ClientParams{
		PrivateSSHKeyPath: gitOpts.PrivateSSHKeyPath,
	})

	logger.Infof("Cloning %s", gitOpts.URL)
	options := git.CloneOptions{
		URL:       gitOpts.URL,
		Branch:    gitOpts.Branch,
		Bootstrap: true,
	}

	cloneDir, err := gitClient.CloneRepoInTmpDir("ops-cli-bootstrap-", options)
	if err != nil {
		return errors.Wrapf(err, "cannot clone repository %s", gitOpts.URL)
	}

	cleanCloneDir := false
	defer func() {
		if cleanCloneDir {
			_ = gitClient.DeleteLocalRepo()
			logger.Debugf("Deleting temporary local clone of %s at %s", gitOpts.URL, cloneDir)
		} else {
			logger.Warningf("You may find the local clone of %s at %s", gitOpts.URL, cloneDir)
		}
	}()

	repoName := "fluxcd"
	repoURL := "https://charts.fluxcd.io"
	var chartName string
	var chartVersion string
	if err := helm.AddRepository(repoName, repoURL); err != nil {
		logger.Fatalf("Failed to add registry %s %s. %s", repoName, repoURL, err)
	}

	chartName = repoName + "/flux"
	chartVersion = "1.1.0"
	if err := helm.PullChartUntarToDir(chartName, chartVersion, cloneDir+"/"+opts.GitFluxPath); err != nil {
		logger.Fatalf("Failed to pull chart %s. %s", chartName, err)
	}

	chartName = repoName + "/helm-operator"
	chartVersion = "0.6.0"
	if err := helm.PullChartUntarToDir(chartName, chartVersion, cloneDir+"/"+opts.GitFluxPath); err != nil {
		logger.Fatalf("Failed to pull chart %s. %s", chartName, err)
	}

	// cleanCloneDir = true

	return nil
}
