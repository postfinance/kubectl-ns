[![Build Status](https://img.shields.io/github/workflow/status/postfinance/kubectl-ns/ci?style=for-the-badge)](https://github.com/postfinance/kubectl-ns/actions)
[![Release](https://img.shields.io/github/release/postfinance/kubectl-ns.svg?style=for-the-badge)](https://github.com/postfinance/kubectl-ns/releases/latest)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge)](/LICENSE.md)
[![Go Report Card](https://img.shields.io/badge/GOREPORT-A%2B-brightgreen.svg?style=for-the-badge)](https://goreportcard.com/report/github.com/postfinance/kubectl-ns)
# kubectl ns plugin
Simple Plugin to display/change the current kube namespace with support for substring matching.

# Build/Installation
## Build from source
go version >= 1.13 with modules support enabled is required to build the plugin from source.
```bash
export GO111MODULES=on # optional if checked out outside of $GOPATH
go build
```
## Installation
Pre-compiled statically linked binaries are available on the [releases page](https://github.com/postfinance/kubectl-ns/releases).

Binary must be placed anywhere in `$PATH` named `kubectl-ns` with execute permissions.
For further information, see the offical documentation on plugins [here](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/).

# Compatibility
Known to work on Windows and Linux. Requires kubectl >= 1.12 (tested with versions >1.12).
Supports the oidc, gcp and azure auth provider for authentication against the k8s api server.

# Examples
For all the examples, assume your cluster has the following namespaces:
```
default
kube-system
kube-public
ingress-nginx
foo
bar
baz
```

## display namespaces
Current namespace is displayed in a different color and last.
```bash
$ kubectl ns
default
kube-system
kube-public
ingress-nginx
foo
bar
baz
```

Substring matching can be used to display namespaces. For example if you are searching for a `kube-` namespace simply type:
```bash
$ kubectl ns kube-
kube-system
kube-public
```

## change current namespace
You can switch the namespace by providing an exact name:
```bash
$ kubectl ns foo
namespace set to "foo"
```

But it's also possible to switch to the `ingress-nginx` namespace by typing a substring (as long as it is a unique name), for example:
```bash
$ kubectl ns ingress
namespace set to "ingress-nginx"
```
