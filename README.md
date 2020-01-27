# ops-cli

## AKSClusterConfig

```yaml
kind: AKSClusterConfig
apiVersion: opscli.io/v1alpha1
metadata:
  name: some-aks-cluster
spec:
  region: west-europe
  version: "1.15.6"
  loadBalancerIP: "10.0.1.2"
  loadBalancerResourceGroup: some-resource-group-ip
  tags:
    managedBy: terraform
```
