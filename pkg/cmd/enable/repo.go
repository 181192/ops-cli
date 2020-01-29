package enable

import (
	api "github.com/181192/ops-cli/pkg/apis/opscli.io/v1alpha1"
	"github.com/181192/ops-cli/pkg/cmd/cmdutils"
	scheme "github.com/181192/ops-cli/pkg/generated/clientset/versioned/scheme"
	"github.com/181192/ops-cli/pkg/git"
	"github.com/181192/ops-cli/pkg/helm"

	"github.com/kr/pretty"
	"github.com/pkg/errors"
	"github.com/rancher/wrangler/pkg/schemes"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func enableRepoCmd() *cobra.Command {
	var gitOpts GitOptions
	var fluxOpts FluxOptions

	cmd := &cobra.Command{
		Use:   "repo",
		Short: "Set up a repo for gitops, installing Flux in the cluster and initializing its manifests.",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			// client, err := kubernetes.NewClient(kubeconfig, configContext)
			// if err != nil {
			// 	return fmt.Errorf("failed to create k8s client: %v", err)
			// }

			// namespace := "monitoring"
			// Create empty cluster config

			scheme.AddToScheme(schemes.All)

			api.NewAKSClusterConfig("", "", api.AKSClusterConfig{})

			// TODO Load from cli or file
			// TODO Create clusterconfig file if not exists in temp repo
			// TODO Add .opsignore file with clusterconfig file

			clusterConfig, err := cmdutils.LoadConfigFromFile(clusterConfigFile)
			if err != nil {
				return err
			}

			logger.Debugf("%# v", pretty.Formatter(clusterConfig))

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
			if err := helm.PullChartUntarToDir(chartName, chartVersion, cloneDir+"/flux-manifests"); err != nil {
				logger.Fatalf("Failed to pull chart %s. %s", chartName, err)
			}

			chartName = repoName + "/helm-operator"
			chartVersion = "0.6.0"
			if err := helm.PullChartUntarToDir(chartName, chartVersion, cloneDir+"/flux-manifests"); err != nil {
				logger.Fatalf("Failed to pull chart %s. %s", chartName, err)
			}

			// cleanCloneDir = true

			return nil
		},
	}

	cmd = configureGitOptions(cmd, &gitOpts)
	cmd = configureFluxOptions(cmd, &fluxOpts)

	return cmd
}
