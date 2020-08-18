package main

import (
	"github.com/181192/ops-cli/cmd"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra/doc"
)

func main() {
	docsCmd := cmd.NewRootCmd()
	docsCmd.DisableAutoGenTag = true
	err := doc.GenMarkdownTree(docsCmd, "./docs")
	if err != nil {
		logger.Fatal(err)
	}
}
