# k3i

> Project k3i is a Kubernetes (k0s) installer for [ignite](https://github.com/weaveworks/ignite) VMs.

It is inspired by [k3d](https://github.com/rancher/k3d)

## Prerequisites

### Install CNI plugins

```bash
export CNI_VERSION=v1.1.1
export ARCH=$([ $(uname -m) = "x86_64" ] && echo amd64 || echo arm64)
sudo mkdir -p /opt/cni/bin
curl -sSL https://github.com/containernetworking/plugins/releases/download/${CNI_VERSION}/cni-plugins-linux-${ARCH}-${CNI_VERSION}.tgz | sudo tar -xz -C /opt/cni/bin
```

### Stop all running containerd services and run a new one

```bash
sudo systemctl stop containerd
sudo killall containerd
sudo killall containerd-shim
sudo killall containerd-shim-runc-v2
sudo containerd -l debug
```

## Installation

### go install
    
```bash
go install github.com/mheers/k3i@latest
```

## Usage

```bash
k3i create test # create a new cluster
k3i list # list all current clusters
k3i kubeconfig test # get the kubeconfig for a cluster
k3i delete test # delete a cluster
```

# TODO
- [ ] add more documentation
- [ ] add more tests
- [ ] support k3s to justify the project name
- [x] implement delete
- [x] implement kubeconfig
- [x] add a command to download assets
- [x] support air-gapped environments

# Alternatives

- https://github.com/innobead/kubefire
- https://github.com/weaveworks/wks-quickstart-firekube
