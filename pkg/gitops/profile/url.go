package profile

import (
	"fmt"

	"github.com/181192/ops-cli/pkg/git"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// RepositoryURL returns the full Git repository URL corresponding to the
// provided "quickstart profile" mnemonic a.k.a. short name. If a valid Git URL
// is provided, this function returns it as-is.
func RepositoryURL(profileArgument string) (string, error) {
	if git.IsGitURL(profileArgument) {
		return profileArgument, nil
	}
	if profile := viper.GetString("profiles." + profileArgument); profile != "" {
		logger.Info(fmt.Sprintf("Using %s profile %s from config %s", profileArgument, profile, viper.ConfigFileUsed()))
		return profile, nil
	}
	if profileArgument == "app-dev" {
		return "https://github.com/weaveworks/eks-quickstart-app-dev", nil
	}
	if profileArgument == "appmesh" {
		return "https://github.com/weaveworks/eks-appmesh-profile", nil
	}
	return "", fmt.Errorf("invalid URL or unknown profile: %s", profileArgument)
}
