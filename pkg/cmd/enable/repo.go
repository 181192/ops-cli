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
	"github.com/181192/ops-cli/pkg/helm"

	"github.com/kr/pretty"
	"github.com/rancher/wrangler/pkg/schemes"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"helm.sh/helm/v3/pkg/cli/values"
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

	if gitOpts.User == "" {
		gitOpts.User = cmd.ClusterConfig.ObjectMeta.Name
	}

	if gitOpts.Email == "" {
		gitOpts.Email = cmd.ClusterConfig.ObjectMeta.Name + "@weave.works"
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

	cloneDir := "./"
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

	fluxManifestsDir := cloneDir + "/" + opts.GitFluxPath

	repoName := "fluxcd"
	repoURL := "https://charts.fluxcd.io"
	if err := helm.AddRepository(repoName, repoURL); err != nil {
		logger.Fatalf("Failed to add registry %s %s. %s", repoName, repoURL, err)
	}

	fluxChartName := repoName + "/flux"
	fluxChartVersion := "1.1.0"
	if err := helm.PullChartUntarToDir(fluxChartName, fluxChartVersion, fluxManifestsDir); err != nil {
		logger.Fatalf("Failed to pull chart %s. %s", fluxChartName, err)
	}

	helmOperatorChartName := repoName + "/helm-operator"
	helmOperatorChartVersion := "0.6.0"
	if err := helm.PullChartUntarToDir(helmOperatorChartName, helmOperatorChartVersion, fluxManifestsDir); err != nil {
		logger.Fatalf("Failed to pull chart %s. %s", helmOperatorChartName, err)
	}

	manifests, err := flux.FillInTemplates(flux.TemplateParameters{})
	if err != nil {
		logger.Fatal("Something went wrong")
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

	writeManifest = func(fileName string, content []byte) error {
		path := filepath.Join(fluxManifestsDir, fileName)
		fmt.Fprintf(os.Stderr, "writing %s\n", path)
		return ioutil.WriteFile(path, content, os.FileMode(0666))
	}

	valueOpts := &values.Options{}

	for fileName, content := range manifests {
		if err := writeManifest(fileName, content); err != nil {
			return fmt.Errorf("cannot output manifest file %s: %s", fileName, err)
		}
	}

	valueOpts.ValueFiles = []string{filepath.Join(fluxManifestsDir, "flux-values.yaml")}
	if err := helm.UpgradeInstallChart("flux", fluxManifestsDir+"/flux", valueOpts); err != nil {
		logger.Fatalf("Failed to pull chart %s. %s", fluxChartName, err)
	}

	valueOpts.ValueFiles = []string{filepath.Join(fluxManifestsDir, "helm-operator-values.yaml")}
	if err := helm.UpgradeInstallChart("helm-operator", fluxManifestsDir+"/helm-operator", valueOpts); err != nil {
		logger.Fatalf("Failed to pull chart %s. %s", fluxChartName, err)
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
