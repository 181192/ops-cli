package main

import (
	"os"

	"github.com/181192/ops-cli/pkg/apis/opscli.io/v1alpha1"

	controllergen "github.com/rancher/wrangler/pkg/controller-gen"
	"github.com/rancher/wrangler/pkg/controller-gen/args"
)

func main() {
	os.Unsetenv("GOPATH")
	controllergen.Run(args.Options{
		OutputPackage: "github.com/181192/ops-cli/pkg/generated",
		Boilerplate:   "hack/boilerplate.go.txt",
		Groups: map[string]args.Group{
			"opscli.io": {
				Types: []interface{}{
					v1alpha1.AKSClusterConfig{},
				},
				GenerateTypes: true,
			},
		},
	})
}
