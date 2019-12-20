package main

import (
	"log"

	"github.com/181192/ops-cli/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	err := doc.GenMarkdownTree(cmd.NewRootCmd(), "./docs")
	if err != nil {
		log.Fatal(err)
	}
}
