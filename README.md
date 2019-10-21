# qserv-operator

A qserv operator for Kubernetes based on [operator-framework](https://github.com/operator-framework). You may be familiar with Operators from the concept’s [introduction in 2016](https://coreos.com/blog/introducing-operators.html). An Operator is a method of packaging, deploying and managing a Kubernetes application.

*operator-sdk version: v0.8.1, commit: 33b3bfe10176f8647f5354516fff29dea42b6342*

[![Build Status](https://travis-ci.org/lsst/qserv-operator.svg?branch=master)](https://travis-ci.org/lsst/qserv-operator)

## Deploy qserv

### Quick start for Ubuntu LTS

```
sudo apt-get update
sudo apt-get install curl docker.io git vim
# then add current user to docker group and restart gnome session
sudo usermod -a -G docker $(id -nu)

WORKDIR="$HOME/src"
mkdir -p "$WORKDIR"

# Create single node k8s cluster with kind
cd "$WORKDIR"
git clone --depth 1 -b "v0.6.0" --single-branch https://github.com/k8s-school/kind-travis-ci
cd kind-travis-ci
./kind/k8s-create.sh -s

cd "$WORKDIR"
git clone  https://github.com/lsst/qserv-operator
cd qserv-operator
./deploy.sh
./wait-operator-ready.sh
kubectl apply -k base
./wait-qserv-ready.sh
./run-integration-tests.sh
```

### Prerequisites

- A valid `KUBECONFIG` and access to a Kubernetes v1.14.2+ cluster
- Dynamic volume provisionning need to be available on the Kubernetes cluster (for example [kind] for or GKE). [kind-travis-ci] provide a one-liner to install [kind] on your workstation.

[kind]:https://kind.sigs.k8s.io/
[kind-travis-ci]:https://github.com/k8s-school/kind-travis-ci

### Deploy qserv-operator and a sample qserv instance 

```sh
# Install qserv-operator
git clone https://github.com/lsst/qserv-operator.git
./deploy.sh

# OPTIONAL: Install a custom qserv instance
# Edit file below to customize this qserv instance
kubectl apply -f deploy/crds/qserv_v1alpha1_qserv_cr.yaml
```

### Connect to the qserv instance

```sh
./run-integration-tests.sh
```

## Build qserv-operator

### Prerequisites

- [git][git_tool]
- [go][go_tool] version v1.12+.
- [docker][docker_tool] version 17.03+.
- [kubectl][kubectl_tool] version v1.11.3+.
- Access to a Kubernetes v1.14.2+ cluster.


[git_tool]:https://git-scm.com/downloads
[go_tool]:https://golang.org/dl/
[docker_tool]:https://docs.docker.com/install/
[kubectl_tool]:https://kubernetes.io/docs/tasks/tools/install-kubectl/

### Build

```sh
git clone https://github.com/kube-incubator/qserv-operator.git
cd qserv-operator
./build-all.sh
```

### Test qserv-operator

```sh
./deploy.sh
./run-multinode-tests.sh
```
