package main

import (
	"github.com/181192/ops-cli/cmd"
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"
)

func main() {
	cmd.Execute()
}
