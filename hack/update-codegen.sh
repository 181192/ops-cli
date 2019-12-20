#!/bin/sh

set -eu

GO111MODULE=on go mod vendor
git clone https://github.com/kubernetes/code-generator.git vendor/k8s.io/code-generator
# check out code-generator at the revision Gopkg.lock had previously referenced
git -C vendor/k8s.io/code-generator checkout aae79feb89bdded3679da91fd8c19b6dfcbdb79a

mkdir -p $HOME/go/src/k8s.io/code-generator
rm -rf $HOME/go/src/k8s.io/code-generator
cp -R vendor/k8s.io/code-generator $HOME/go/src/k8s.io/

# ROOT_PACKAGE :: the package that is the target for code generation
ROOT_PACKAGE=github.com/181192/ops-cli
# CUSTOM_RESOURCE_NAME :: the name of the custom resource that we're generating client code for
CUSTOM_RESOURCE_NAME=aks
# CUSTOM_RESOURCE_VERSION :: the version of the resource
CUSTOM_RESOURCE_VERSION=v1alpha1

bindir=$( cd "${0%/*}" && pwd )
rootdir=$( cd "$bindir"/.. && pwd )

# run the code-generator entrypoint script
"$rootdir"/vendor/k8s.io/code-generator/generate-groups.sh deepcopy "$ROOT_PACKAGE/pkg/gen/client" "$ROOT_PACKAGE/pkg/gen/apis" "$CUSTOM_RESOURCE_NAME:$CUSTOM_RESOURCE_VERSION" $@

cp -r $HOME/go/src/github.com/181192/ops-cli/pkg/gen "$rootdir"/pkg
rm -rf $HOME/go/src/github.com/181192/ops-cli
