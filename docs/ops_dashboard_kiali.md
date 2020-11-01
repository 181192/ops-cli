## ops dashboard kiali



### Synopsis



```
ops dashboard kiali [flags]
```

### Options

```
  -h, --help                    help for kiali
      --kube-context string     Name of the kubeconfig context to use
      --kubeconfig string       Absolute path of the kubeconfig file to be used
  -l, --label-selector string   Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)
  -n, --namespace string        Name of the namespace to use. Defaults to the application default namespace. (default "istio-system")
  -p, --port int                Target port to forward to
```

### Options inherited from parent commands

```
      --log-level string   Log level (debug, info, warn, error, fatal, panic) (default "info")
```

### SEE ALSO

* [ops dashboard](ops_dashboard.md)	 - Dashboards  - alertmanager
  - grafana
  - jaeger
  - kiali
  - prometheus

