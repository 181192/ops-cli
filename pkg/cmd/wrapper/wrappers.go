package wrapper

import "github.com/181192/ops-cli/pkg/wrapper"

// Wrapper represents details of a wrapper command
type Wrapper struct {
	Name            string
	Description     string
	LongDescription string
	Executable      string
}

// MakeWrappers returns list of all wrapper commands
func MakeWrappers() []Wrapper {
	wrappers := []Wrapper{}

	wrappers = append(wrappers, Wrapper{
		Name:        "helm",
		Description: "A kubernetes package manager",
		Executable:  wrapper.HelmBinary,
	})

	wrappers = append(wrappers, Wrapper{
		Name:        "helmfile",
		Description: "Deploy Kubernetes Helm Charts",
		LongDescription: `Helmfile is a declarative spec for deploying helm charts. It lets you...

		- Keep a directory of chart value files and maintain changes in version control.
		- Apply CI/CD to configuration changes.
		- Periodically sync to avoid skew in environments.`,
		Executable: wrapper.HelmfileBinary,
	})

	wrappers = append(wrappers, Wrapper{
		Name:        "kubectl",
		Description: "kubectl controls the Kubernetes cluster manager",
		Executable:  wrapper.KubectlBinary,
	})

	wrappers = append(wrappers, Wrapper{
		Name:            "terraform",
		Description:     "Terraform IaC tool",
		LongDescription: "Terraform IaC tool",
		Executable:      wrapper.TerraformBinary,
	})

	return wrappers
}
