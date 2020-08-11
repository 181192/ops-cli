package wrapper

import (
	"os"

	"github.com/181192/ops-cli/pkg/util"
)

var (
	// HelmBinary path to helm binary
	HelmBinary = getBinaryPath("helm")

	// HelmfileBinary path to helmfile binary
	HelmfileBinary = getBinaryPath("helmfile")

	// TerraformBinary path to terraform binary
	TerraformBinary = getBinaryPath("terraform")

	// KubectlBinary path to kubectl binary
	KubectlBinary = getBinaryPath("kubectl")
)

func getBinaryPath(binary string) string {
	return util.GetConfigDirectory() + string(os.PathSeparator) + "bin" + string(os.PathSeparator) + binary + util.GetWinExtension()
}
