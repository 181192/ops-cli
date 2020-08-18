package cmdutils

import (
	"github.com/spf13/pflag"
)

const (
	kubeConfig    = "kubeconfig"
	kubeContext   = "kube-context"
	kubeNamespace = "namespace"
)

// KubernetesOpts common kubernetes opts
type KubernetesOpts struct {
	KubeConfig  string
	KubeContext string
	Namespace   string
}

// AddCommonFlagsForKubernetes configures the flags required to interact with Kubernetes cluster
func AddCommonFlagsForKubernetes(fs *pflag.FlagSet, opts *KubernetesOpts, defaultNamespace string) {
	fs.StringVar(&opts.KubeContext, kubeContext, "", "Name of the kubeconfig context to use")
	fs.StringVar(&opts.KubeConfig, kubeConfig, "", "Absolute path of the kubeconfig file to be used")
	fs.StringVarP(&opts.Namespace, kubeNamespace, "n", defaultNamespace, "Name of the namespace to use. Defaults to the application default namespace.")
}
