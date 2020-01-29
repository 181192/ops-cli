package enable

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/181192/ops-cli/pkg/git"
	"github.com/181192/ops-cli/pkg/helm"
	"github.com/kr/pretty"

	api "github.com/181192/ops-cli/pkg/apis/opscli.io/v1alpha1"
	scheme "github.com/181192/ops-cli/pkg/generated/clientset/versioned/scheme"

	"github.com/pkg/errors"
	"github.com/rancher/wrangler/pkg/schemes"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
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

			clusterConfig, err := LoadConfigFromFile(clusterConfigFile)
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

// LoadConfigFromFile loads ClusterConfig from configFile
func LoadConfigFromFile(configFile string) (*api.AKSClusterConfig, error) {
	data, err := readConfig(configFile)
	if err != nil {
		return nil, errors.Wrapf(err, "reading config file %q", configFile)
	}

	if err := yaml.UnmarshalStrict(data, &api.AKSClusterConfig{}); err != nil {
		return nil, errors.Wrapf(err, "loading config file %q", configFile)
	}

	obj, err := runtime.Decode(scheme.Codecs.UniversalDeserializer(), data)
	if err != nil {
		return nil, errors.Wrapf(err, "loading config file %q", configFile)
	}

	cfg, ok := obj.(*api.AKSClusterConfig)
	if !ok {
		return nil, fmt.Errorf("expected to decode object of type %T; got %T", &api.AKSClusterConfig{}, cfg)
	}
	return cfg, nil
}

func readConfig(configFile string) ([]byte, error) {
	if configFile == "-" {
		return ioutil.ReadAll(os.Stdin)
	}
	return ioutil.ReadFile(configFile)
}
