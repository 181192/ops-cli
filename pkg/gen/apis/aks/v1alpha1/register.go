package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"

	aks "github.com/181192/ops-cli/pkg/gen/apis/aks"
)

// Conventional Kubernetes API contants
const (
	CurrentGroupVersion = "v1alpha1"
	ClusterConfigKind   = "ClusterConfig"
)

// Conventional Kubernetes API variables
var (
	SchemeGroupVersion = schema.GroupVersion{
		Group:   aks.GroupName,
		Version: CurrentGroupVersion,
	}

	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

// Kind takes an unqualified kind and returns back a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

// Register our API with the scheme
func Register() error {
	return AddToScheme(scheme.Scheme)
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&ClusterConfig{},
		&ClusterConfigList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
