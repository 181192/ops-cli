## ops enable repo

Set up a repo for gitops, installing Flux in the cluster and initializing its manifests

```
ops enable repo [flags]
```

### Options

```
      --acr-registry                         Enable ACR authentication (requires deployment in AKS) (default true)
  -c, --cluster string                       Cluster name
  -f, --config-file string                   Load configuration from a file (or stdin if set to '-')
      --flux-chart-version string            Chart version of Flux (default latest)
      --garbage-collection                   Enable garbage collection (default true)
      --git-branch string                    Git branch to be used for GitOps (default "master")
      --git-email string                     Email to use as Git committer
      --git-flux-subdir string               Directory within the Git repository where to commit the Flux manifests (default "manifests-flux/")
      --git-label string                     Git label to keep track of Flux's sync progress; this is equivalent to overriding --git-sync-tag and --git-notes-ref in Flux
      --git-paths strings                    Relative paths within the Git repo for Flux to locate Kubernetes manifests
      --git-private-ssh-key-path string      Optional path to the private SSH key to use with Git, e.g. ~/.ssh/id_rsa
      --git-url string                       SSH URL of the Git repository to be used for GitOps, e.g. git@github.com:<github_org>/<repo_name>
      --git-user string                      Username to use as Git committer
      --helm-operator-chart-version string   Chart version of Helm Operator (default latest)
      --helm-versions strings                Versions of Helm to enable (default [v3])
  -h, --help                                 help for repo
      --kube-context string                  Name of the kubeconfig context to use
      --kubeconfig string                    Absolute path of the kubeconfig file to be used
  -l, --location string                      Cluster location (default "westeurope")
      --manifest-generation                  Enable manifest generation (default true)
  -n, --namespace string                     Name of the namespace to use. Defaults to the application default namespace. (default "flux-system")
      --override-values                      Override values files
      --skip-install                         Skip installing Flux to cluster
      --with-helm                            Install the Helm Operator (default true)
```

### Options inherited from parent commands

```
      --config string      Config file (default is /home/k/.ops/ops.[yaml|json|toml|properties])
      --log-level string   Log level (debug, info, warn, error, fatal, panic) (default "info")
```

### SEE ALSO

* [ops enable](ops_enable.md)	 - Enable resource(s)

