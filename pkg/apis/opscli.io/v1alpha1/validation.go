package v1alpha1

// ValidateAKSClusterConfig validates AKSClusterConfig
func ValidateAKSClusterConfig(cfg *AKSClusterConfig) error {
	// TODO
	return nil
}

// SupportedVersions are the supported versions of Kubernetes
func SupportedVersions() []string {
	return []string{
		Version1_13,
		Version1_14,
		Version1_15,
		Version1_16,
	}
}
