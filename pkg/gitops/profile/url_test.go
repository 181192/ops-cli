package profile_test

import (
	"testing"

	"github.com/181192/ops-cli/pkg/gitops/profile"
	"github.com/181192/ops-cli/pkg/testutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSuite(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("profile", func() {
	Describe("RepositoryURL", func() {
		It("returns Git URLs as-is", func() {
			url, err := profile.RepositoryURL("https://github.com/eksctl-bot/my-gitops-repo")
			Expect(err).To(Not(HaveOccurred()))
			Expect(url).To(Equal("https://github.com/eksctl-bot/my-gitops-repo"))
		})

		It("returns full Git URLs for supported mnemonics", func() {
			mnemonicToURLs := []struct {
				mnemonic string
				url      string
			}{
				{mnemonic: "app-dev", url: "https://github.com/weaveworks/eks-quickstart-app-dev"},
				{mnemonic: "appmesh", url: "https://github.com/weaveworks/eks-appmesh-profile"},
			}
			for _, mnemonicToURL := range mnemonicToURLs {
				url, err := profile.RepositoryURL(mnemonicToURL.mnemonic)
				Expect(err).To(Not(HaveOccurred()))
				Expect(url).To(Equal(mnemonicToURL.url))
			}
		})

		It("returns an error otherwise", func() {
			url, err := profile.RepositoryURL("foo")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("invalid URL or unknown profile: foo"))
			Expect(url).To(Equal(""))
		})
	})
})
