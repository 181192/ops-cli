package ops

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/yaml"

	api "github.com/181192/ops-cli/pkg/apis/opscli.io/v1alpha1"
)

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
