module github.com/181192/ops-cli

go 1.15

require (
	github.com/181192/ops-cli/internal v0.1.16
	github.com/cheggaaa/pb v1.0.27
	github.com/gofrs/flock v0.8.1
	github.com/hashicorp/go-getter v1.4.1
	github.com/hashicorp/go-version v1.4.0
	github.com/kr/pretty v0.3.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.10.1
	github.com/pkg/errors v0.9.1
	github.com/rancher/wrangler v0.6.1
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/afero v1.8.2
	github.com/spf13/cobra v1.4.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.2
	github.com/stretchr/testify v1.7.1
	github.com/tidwall/gjson v1.6.0
	github.com/whilp/git-urls v0.0.0-20191001220047-6db9661140c0
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1
	gopkg.in/yaml.v2 v2.4.0
	helm.sh/helm/v3 v3.2.1
	k8s.io/api v0.18.0
	k8s.io/apimachinery v0.18.0
	k8s.io/client-go v0.18.0
	k8s.io/kops v1.11.0
	rsc.io/letsencrypt v0.0.3 // indirect
	sigs.k8s.io/yaml v1.2.0
)

replace github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309

replace k8s.io/kubectl => k8s.io/kubectl v0.18.0

replace github.com/Azure/go-autorest/autorest => github.com/Azure/go-autorest/autorest v0.11.15

replace github.com/Azure/go-autorest/autorest/adal => github.com/Azure/go-autorest/autorest/adal v0.9.10

replace github.com/181192/ops-cli/internal => ./internal
