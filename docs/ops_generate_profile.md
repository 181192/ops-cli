## ops generate profile

Generate a gitops profile

### Synopsis

Generate a gitops profile

```
ops generate profile [flags]
```

### Options

```
  -c, --cluster string                    Cluster name
  -f, --config-file string                Load configuration from a file (or stdin if set to '-')
      --git-branch string                 Git branch (default "master")
      --git-private-ssh-key-path string   Optional path to the private SSH key to use with Git, e.g. ~/.ssh/id_rsa
      --git-url string                    URL for the quickstart base repository
  -h, --help                              help for profile
  -l, --location string                   Cluster location (default "westeurope")
      --profile-path string               Path to generate the profile in (default "./")
```

### Options inherited from parent commands

```
      --config string      Config file (default is /home/k/.ops/ops.[yaml|json|toml|properties])
      --log-level string   Log level (debug, info, warn, error, fatal, panic) (default "info")
```

### SEE ALSO

* [ops generate](ops_generate.md)	 - Generate resource(s)

