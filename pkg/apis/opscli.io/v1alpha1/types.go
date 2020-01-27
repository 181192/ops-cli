package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AKSClusterConfig represent a cluster
type AKSClusterConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              *AKSClusterConfigSpec `json:"spec,omitempty"`
}

// AKSClusterConfigSpec is what identifies a cluster
type AKSClusterConfigSpec struct {
	Location                  string `json:"location"`
	LoadBalancerIP            string `json:"loadBalancerIp"`
	LoadBalancerResourceGroup string `json:"loadBalancerResourceGroup"`
	// +optional
	Version string `json:"version,omitempty"`
	// +optional
	Tags map[string]string `json:"tags,omitempty"`
}
