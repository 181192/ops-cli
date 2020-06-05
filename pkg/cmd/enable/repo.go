package enable

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	api "github.com/181192/ops-cli/pkg/apis/opscli.io/v1alpha1"
	"github.com/181192/ops-cli/pkg/cmd/cmdutils"
	"github.com/181192/ops-cli/pkg/flux"
	scheme "github.com/181192/ops-cli/pkg/generated/clientset/versioned/scheme"
	"github.com/181192/ops-cli/pkg/git/gitconfig"
	"github.com/181192/ops-cli/pkg/helm"
	"github.com/181192/ops-cli/pkg/util/file"
	"helm.sh/helm/v3/pkg/cli/values"

	"github.com/kr/pretty"
	"github.com/rancher/wrangler/pkg/schemes"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	repoName              = "fluxcd"
	repoURL               = "https://charts.fluxcd.io"
	fluxChartName         = repoName + "/flux"
	helmOperatorChartName = repoName + "/helm-operator"
)

func enableRepoCmd(cmd *cmdutils.Cmd) {
	var opts cmdutils.InstallOpts

	cmd.ClusterConfig = api.DefaultClusterConfig()
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
	gitOpts := &opts.GitOptions

	if err := cmdutils.NewMetadataLoader(cmd).Load(); err != nil {
		return err
	}

	clusterName := cmd.ClusterConfig.ObjectMeta.Name

	if gitOpts.URL == "" {
		url, err := gitconfig.OriginURL()
		if err != nil {
			return err
		}
		gitOpts.URL = url
	}

	if gitOpts.User == "" {
		gitOpts.User = clusterName
	}

	if gitOpts.Email == "" {
		gitOpts.Email = clusterName + "@weave.works"
	}

	if opts.GitLabel == "" {
		opts.GitLabel = clusterName + "-sync"
	}

	if err := cmdutils.ValidateGitOptions(gitOpts); err != nil {
		return err
	}

	if err := cmdutils.NewGitOpsConfigLoader(cmd).Load(); err != nil {
		return err
	}

	logger.Debugf("%# v", pretty.Formatter(cmd.ClusterConfig))
	logger.Debugf("%# v", pretty.Formatter(opts))

	// gitClient := git.NewGitClient(git.ClientParams{
	// 	PrivateSSHKeyPath: gitOpts.PrivateSSHKeyPath,
	// })

	// logger.Infof("Cloning %s", gitOpts.URL)
	// options := git.CloneOptions{
	// 	URL:       gitOpts.URL,
	// 	Branch:    gitOpts.Branch,
	// 	Bootstrap: true,
	// }

	cloneDir := "." + string(os.PathSeparator)
	// cloneDir, err := gitClient.CloneRepoInTmpDir("ops-cli-bootstrap-", options)
	// if err != nil {
	// 	return errors.Wrapf(err, "cannot clone repository %s", gitOpts.URL)
	// }

	// cleanCloneDir := false
	// defer func() {
	// 	if cleanCloneDir {
	// 		_ = gitClient.DeleteLocalRepo()
	// 		logger.Debugf("Deleting temporary local clone of %s at %s", gitOpts.URL, cloneDir)
	// 	} else {
	// 		logger.Warningf("You may find the local clone of %s at %s", gitOpts.URL, cloneDir)
	// 	}
	// }()

	fluxManifestsDir := cloneDir + string(os.PathSeparator) + opts.GitFluxPath

	if err := helm.AddRepository(repoName, repoURL); err != nil {
		logger.Fatalf("Failed to add registry %s %s. %s", repoName, repoURL, err)
	}

	if err := helm.PullChartUntarToDir(fluxChartName, opts.FluxChartVersion, fluxManifestsDir); err != nil {
		logger.Fatalf("Failed to pull chart %s. %s", fluxChartName, err)
	}

	if err := helm.PullChartUntarToDir(helmOperatorChartName, opts.HelmOperatorChartVersion, fluxManifestsDir); err != nil {
		logger.Fatalf("Failed to pull chart %s. %s", helmOperatorChartName, err)
	}

	templateParameters := &flux.TemplateParameters{
		GitURL:             opts.GitOptions.URL,
		GitBranch:          opts.GitOptions.Branch,
		GitUser:            opts.GitOptions.User,
		GitEmail:           opts.GitOptions.Email,
		GitLabel:           opts.GitLabel,
		GitPaths:           opts.GitPaths,
		AcrRegistry:        opts.AcrRegistry,
		GarbageCollection:  opts.GarbageCollection,
		ManifestGeneration: opts.ManifestGeneration,
		HelmVersions:       opts.HelmVersions,
	}
	manifests, err := flux.FillInTemplates(templateParameters)
	if err != nil {
		logger.Fatalf("Failed to template install values %s", err)
	}

	writeManifest := func(fileName string, content []byte) error {
		_, err := os.Stdout.Write(content)
		return err
	}

	info, err := os.Stat(fluxManifestsDir)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", fluxManifestsDir)
	}

	writeManifest = func(path string, content []byte) error {
		logger.Infof("Writing %s", path)
		return ioutil.WriteFile(path, content, os.FileMode(0666))
	}

	for fileName, content := range manifests {
		fileName = clusterName + "-" + fileName
		path := filepath.Join(fluxManifestsDir, fileName)
		if !file.Exists(path) || opts.OverrideValues {
			if err := writeManifest(path, content); err != nil {
				return fmt.Errorf("cannot output manifest file %s: %s", path, err)
			}
		} else {
			logger.Warnf("%s exists, skip creating...", path)
		}
	}

	if opts.SkipInstall {
		logger.Warning("Skip installing Flux to cluster...")
		return nil
	}

	valueOpts := &values.Options{}

	fluxValues := clusterName + "-flux-values.yaml"
	fluxChartPath := fluxManifestsDir + string(os.PathSeparator) + "flux"
	valueOpts.ValueFiles = []string{filepath.Join(fluxManifestsDir, fluxValues)}
	if err := helm.UpgradeInstallChart("flux", fluxChartPath, valueOpts, opts); err != nil {
		logger.Fatalf("Failed to install chart %s. %s", fluxChartName, err)
	}

	helmOperatorValues := clusterName + "-helm-operator-values.yaml"
	helmOperatorChartPath := fluxManifestsDir + string(os.PathSeparator) + "helm-operator"
	valueOpts.ValueFiles = []string{filepath.Join(fluxManifestsDir, helmOperatorValues)}
	if err := helm.UpgradeInstallChart("helm-operator", helmOperatorChartPath, valueOpts, opts); err != nil {
		logger.Fatalf("Failed to install chart %s. %s", helmOperatorChartName, err)
	}

	// Git add, commit and push flux manifests files in the user's repo
	// if err = gitClient.Add("."); err != nil {
	// 	return err
	// }

	// userGitUser, err := gitconfig.Username()
	// if err != nil {
	// 	return err
	// }

	// userGitEmail, err := gitconfig.Email()
	// if err != nil {
	// 	return err
	// }

	// if err = gitClient.Commit("Add flux manifests", userGitUser, userGitEmail); err != nil {
	// 	return err
	// }

	// if err = gitClient.Push(); err != nil {
	// 	return err
	// }

	// cleanCloneDir = true

	return nil
}
