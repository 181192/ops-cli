# ops-cli

## Installation

```
# Darwin/macOS
curl -sSL https://github.com/181192/ops-cli/releases/latest/download/ops_cli_darwin_amd64 -o ops \
  && chmod +x ops \
  && sudo mv ops /usr/local/bin/ops

# Linux
curl -sSL https://github.com/181192/ops-cli/releases/latest/download/ops_cli_linux_amd64 -o ops \
  && chmod +x ops \
  && sudo mv ops /usr/local/bin/ops

# Windows (run terminal as admin)
Invoke-WebRequest -Uri https://github.com/181192/ops-cli/releases/latest/download/ops_cli_windows_amd64.exe -OutFile ops.exe
```

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

## Configuration file

Create a `ops.yaml` file in `$HOME/.ops` or current directory

### Custom dashboards

Example config:

```yaml
dashboards:
  - name: cadvisor
    namespace: cadvisor
    port: 8080
    labelSelector: app=cadvisor
    url: /metrics
```

Will open the browser at http://localhost:8080/metrics

Can also override url og provive a url if not set in config with:

```
ops dashboard grafana --url /explore

# or
ops d grafana -u /explore
```

OR set default grafana url to always be /explore by overriding config

```yaml
dashboards:
  - name: grafana
    namespace: monitoring
    port: 3000
    labelSelector: app.kubernetes.io/name=grafana
    url: /explore
```

## Default profiles

Adding profiles configuration to the config file makes it easier to use ops-cli when updating a repo.

```yaml
profiles:
  default: git@github.com:181192/kustomize-manifests.git
  example: git@github.com:181192/some-other-manifests.git
```

The ops enable profile command simplifies from:

```
ops enable profile git@github.com:181192/kustomize-manifests.git --manifest-only

# or
ops enable profile --name git@github.com:181192/kustomize-manifests.git --manifest-only
```

To using the default profile:

```
ops enable profile --manifest-only
```

Or using the example profile:

```
ops enable profile example --manifest-only

# or
ops enable profile --name example --manifest-only
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
