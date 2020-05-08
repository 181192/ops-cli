package wrapper

import "github.com/181192/ops-cli/pkg/util"

var (
	HelmBinary      = util.GetConfigDirectory() + "/bin/helm"
	HelmfileBinary  = util.GetConfigDirectory() + "/bin/helmfile"
	TerraformBinary = util.GetConfigDirectory() + "/bin/terraform"
	KubectlBinary   = util.GetConfigDirectory() + "/bin/kubectl"
)
