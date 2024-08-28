# Kubectl Plugin: Pod Resource Summary
**🚧 Work In Progress 🚧**

### Purpose:

To provide a quick and concise summary of the resource usage (CPU, memory) and status (running, pending, etc.) of all pods within a namespace or across the entire cluster. This helps developers and SREs get a quick snapshot of the cluster’s health.

### Example
```shell
NAMESPACE             PODS  RUNNING  PENDING   FAILED   CPU(m) MEMORY(Mi)
-------------------------------------------------------------------------------
kube-system              8        8        0        0      950        290
litmus                   1        0        0        1      100        128
local-path-storage       1        1        0        0        0          0
```