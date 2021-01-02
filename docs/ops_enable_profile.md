## ops enable profile

Enable and deploy the components from the selected profile

### Synopsis

Enable and deploy the components from the selected profile

```
ops enable profile [flags]
```

### Options

```
  -c, --cluster string                    Cluster name
  -f, --config-file string                Load configuration from a file (or stdin if set to '-')
      --git-branch string                 Git branch to be used for GitOps (default "master")
      --git-email string                  Email to use as Git committer
      --git-private-ssh-key-path string   Optional path to the private SSH key to use with Git, e.g. ~/.ssh/id_rsa
      --git-url string                    SSH URL of the Git repository to be used for GitOps, e.g. git@github.com:<github_org>/<repo_name>
      --git-user string                   Username to use as Git committer
  -h, --help                              help for profile
      --load-balacer-ip string            Loadbalancer IP
      --load-balacer-ip-rg string         Loadbalancer IP resource group
  -l, --location string                   Cluster location (default "westeurope")
      --manifest-only                     Only update manifests directory, ignore profile.
      --name string                       Name or URL of the profile. For example, app-dev.
      --overlay string                    Name of the overlay profile. For example nginx,linkerd or istio. (default "nginx")
      --revision string                   Revision of the profile. (default "master")
```

### Options inherited from parent commands

```
      --config string      Config file (default is /home/k/.ops/ops.[yaml|json|toml|properties])
      --log-level string   Log level (debug, info, warn, error, fatal, panic) (default "info")
```

### SEE ALSO

* [ops enable](ops_enable.md)	 - Enable resource(s)

