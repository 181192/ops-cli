package enable

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"

	"sigs.k8s.io/yaml"

	"github.com/pkg/errors"

	"io/ioutil"
	"os"

	api "github.com/181192/ops-cli/pkg/apis/opscli.io/v1alpha1"
	"github.com/181192/ops-cli/pkg/git"

	"github.com/spf13/cobra"
)

func enableRepoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repo",
		Short: "Set up a repo for gitops, installing Flux in the cluster and initializing its manifests.",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			// client, err := kubernetes.NewClient(kubeconfig, configContext)
			// if err != nil {
			// 	return fmt.Errorf("failed to create k8s client: %v", err)
			// }

			// namespace := "monitoring"
			// Create empty cluster config
			api.NewAKSClusterConfig("", "", api.AKSClusterConfig{})

			// Parse args / config file

			// print config

			return nil
		},
	}

	return cmd
}

type InstallOpts struct {
	GitOptions  git.Options
	GitPaths    []string
	GitLabel    string
	GitFluxPath string
	Namespace   string
	Timeout     time.Duration
	WithHelm    bool
}

func configureRepositoryCmd(cmd *cobra.Command) {

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
