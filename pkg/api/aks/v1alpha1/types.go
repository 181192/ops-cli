package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Version1_15 represents Kubernetes version 1.15.x
	Version1_15 = "1.15"

	// DefaultVersion represents default Kubernetes version supported by AKS
	DefaultVersion = Version1_15

	// LatestVersion represents latest Kubernetes version supported by AKS
	LatestVersion = Version1_15
)

// ClusterMeta is what identifies a cluster
type ClusterMeta struct {
	Name   string `json:"name"`
	Region string `json:"region"`
	// +optional
	Version string `json:"version,omitempty"`
	// +optional
	Tags map[string]string `json:"tags,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterConfig is a simple config, to be replaced with Cluster API
type ClusterConfig struct {
	metav1.TypeMeta `json:",inline"`

	Metadata *ClusterMeta `json:"metadata"`
}

// ClusterConfigTypeMeta constructs TypeMeta for ClusterConfig
func ClusterConfigTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       ClusterConfigKind,
		APIVersion: SchemeGroupVersion.String(),
	}
}

// NewClusterConfig creates new config for a cluster;
// it doesn't include initial nodegroup, so user must
// call NewNodeGroup to create one
func NewClusterConfig() *ClusterConfig {
	cfg := &ClusterConfig{
		TypeMeta: ClusterConfigTypeMeta(),
		Metadata: &ClusterMeta{
			Version: DefaultVersion,
		},
	}

	return cfg
}
