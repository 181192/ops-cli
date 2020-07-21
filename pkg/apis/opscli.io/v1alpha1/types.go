package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Version1_13 represents Kubernetes version 1.13.x
	Version1_13 = "1.13"

	// Version1_14 represents Kubernetes version 1.14.x
	Version1_14 = "1.14"

	// Version1_15 represents Kubernetes version 1.15.x
	Version1_15 = "1.15"

	// Version1_16 represents Kubernetes version 1.16.x
	Version1_16 = "1.16"

	// DefaultVersion represents default Kubernetes version supported
	DefaultVersion = Version1_15

	// LatestVersion represents latest Kubernetes version supported
	LatestVersion = Version1_16
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterConfig represent a cluster
type ClusterConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              *ClusterConfigSpec `json:"spec,omitempty"`
}

// ClusterConfigSpec is what identifies a cluster
type ClusterConfigSpec struct {
	Location                  string `json:"location"`
	LoadBalancerIP            string `json:"loadBalancerIP"`
	LoadBalancerResourceGroup string `json:"loadBalancerResourceGroup"`
	// +optional
	Version string `json:"version,omitempty"`
	// +optional
	Tags map[string]string `json:"tags,omitempty"`
}

// DefaultClusterConfig creates new config for cluster
func DefaultClusterConfig() *ClusterConfig {
	return &ClusterConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterConfig",
			APIVersion: "v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       &ClusterConfigSpec{},
	}
}
