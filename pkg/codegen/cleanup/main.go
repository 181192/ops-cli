package main

import (
	"os"

	"github.com/rancher/wrangler/pkg/cleanup"
	logger "github.com/sirupsen/logrus"
)

func main() {
	if err := cleanup.Cleanup("./pkg/apis"); err != nil {
		logger.Fatal(err)
	}
	if err := os.RemoveAll("./pkg/generated"); err != nil {
		logger.Fatal(err)
	}
}
