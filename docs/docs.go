package main

import (
	"github.com/181192/ops-cli/cmd"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra/doc"
)

func main() {
	err := doc.GenMarkdownTree(cmd.NewRootCmd(), "./docs")
	if err != nil {
		logger.Fatal(err)
	}
}
