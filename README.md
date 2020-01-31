# ops-cli

## Example workflow
1. Create a `cluster-config.yaml`
```
cat <<EOF > cluster-config.yaml
kind: ClusterConfig
apiVersion: opscli.io/v1alpha1
metadata:
  name: some-aks-cluster
spec:
  location: westeurope
  version: "1.15.6"
  loadBalancerIP: "10.0.1.2"
  loadBalancerResourceGroup: some-resource-group-ip
EOF
```

2. Create a new repository on GitHub/Bitbucket...

3. Install and enable Flux & Helm-Operator in your cluster and save the manifests in the git repository
```
ops enable repo \
  --git-url git@github.com:181192/empty-sample.git \
  -f cluster-config.yaml
```

4. Enable a GitOps profile and deploy to cluster
```
ops enable profile \
  --name git@bitbucket.org:181192/kustomize-manifests.git \
  --git-url git@github.com:181192/empty-sample.git \
  -f cluster-config.yaml
```

## Building from sources

Building the project
```
git clone git@github.com:181192/ops-cli.git && cd ops-cli
make deps
make build
```

Generate docs
```
make docs
```

Generate CRD API's
```
make generate-cleanup
make generate
```
