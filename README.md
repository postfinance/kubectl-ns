# kubectl ns plugin
Simple Plugin to display/change the current kube namespace.  

# Build/Installation
## Build from source
go version >= 1.11 with modules support enabled is required to build the plugin from source
```bash
export GO111MODULES=on
go build
```
## Installation
Binary must be placed anywhere in `$PATH` named `kubectl-ns` with execute permissions.  
For further information, see the offical documentation on plugins [here](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/).

# Compatibility
Known to work on Windows and Linux. Requires kubectl >= 1.12 (tested with 1.12).
Supports the oidc, gcp and azure auth provider for authentication against the k8s api server.

# Examples
## display namespaces
Current namespace is displayed in a different color.
```bash
$ kubectl ns
foo
bar
baz
```

## change current namespace
```bash
$ kubectl ns foo
namespace set to "foo"
```
