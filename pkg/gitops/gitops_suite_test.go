package gitops

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/181192/ops-cli/pkg/testutils"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	testutils.RegisterAndRun(t)
}
