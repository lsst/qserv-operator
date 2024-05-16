###########
Quick start
###########

Prerequisites
=============

* An Ubuntu LTS workstation
* Internet access without proxy
* `sudo` access

Install dependencies and add user to `docker` group
---------------------------------------------------

.. code:: bash

    sudo apt-get update
    sudo apt-get install curl docker.io git vim
    sudo usermod -a -G docker $(id -nu)

.. warning::

    Restart session in order to take in account add to `docker` group.

Create a single node k8s cluster
--------------------------------

Option #1: kind
^^^^^^^^^^^^^^^

`kind <https://kind.sigs.k8s.io/>`__ is a tool for running local Kubernetes clusters using Docker container “nodes”.
kind was primarily designed for testing Kubernetes itself, but may be used for local development or CI.
Script below uses a `simple install script for kind <https://github.com/k8s-school/kind-helper>`__ provided by `K8s-school <https://k8s-school.fr>`__.

.. code:: bash

    WORKDIR="$HOME/src"
    mkdir -p "$WORKDIR"

    cd "$WORKDIR"
    git clone --depth 1 -b "k8s-v1.20.2" --single-branch https://github.com/k8s-school/kind-helper
    cd kind-helper
    ./kind/k8s-create.sh -s

Option #2: k3s
^^^^^^^^^^^^^^

`k3s <https://k3s.io/>`__ is the certified Kubernetes distribution built for IoT & Edge computing. It may be used for local development or CI.

.. code:: bash

    curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC=“--docker --write-kubeconfig-mode 644” sh -
    export KUBECONFIG=/etc/rancher/k3s/k3s.yaml

Deploy Qserv in four lines
===========================

This procedure is recommended for development platform only.

`golang <https://go.dev/>`__ and `git-lfs <https://git-lfs.com>`__ are pre-requisites.

.. code:: bash

    RELEASE="2024.5.1-rc4"
    git clone --depth 1 --single-branch -b "$RELEASE" https://github.com/lsst/qserv-operator
    cd qserv-operator
    # Install pre-requisites
    ./prereq-install.sh
    # Deploy Qserv operator
    kubectl apply -f manifests/operator.yaml
    # Deploy Qserv
    kubectl apply -k manifests/base
    # Run integration tests
    ./tests/tools/wait-qserv-ready.sh
    ./tests/e2e/integration.sh



