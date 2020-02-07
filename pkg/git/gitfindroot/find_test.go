package gitfindroot

import (
	"testing"

	logger "github.com/sirupsen/logrus"
)

func TestRootIsFound(t *testing.T) {
	response, err := Repo()
	if err != nil {
		logger.Fatalf("Error: %s", err.Error())
	}

	expectation := "ops-cli"

	if response.Name != expectation {
		t.Errorf("The response '%s' didn't match the expectaton '%s'", response.Name, expectation)
	}
}
