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
Script below uses a `simple install script for kind <https://github.com/k8s-school/kind-travis-ci>`__ provided by `K8s-school <https://k8s-school.fr>`__.

.. code:: bash

    WORKDIR="$HOME/src"
    mkdir -p "$WORKDIR"

    cd "$WORKDIR"
    git clone --depth 1 -b "v0.6.0" --single-branch https://github.com/k8s-school/kind-travis-ci
    cd kind-travis-ci
    ./kind/k8s-create.sh -s

Option #2: k3s
^^^^^^^^^^^^^^

`k3s <https://k3s.io/>`__ is the certified Kubernetes distribution built for IoT & Edge computing. It may be used for local development or CI.

.. code:: bash

    curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC=“--docker --write-kubeconfig-mode 644” sh -
    export KUBECONFIG=/etc/rancher/k3s/k3s.yaml

Install Qserv in two lines
==========================

.. code:: bash

    curl -fsSL https://raw.githubusercontent.com/lsst/qserv-operator/master/deploy/qserv.sh | bash -s
    kubectl apply -k https://github.com/lsst/qserv-operator/base

Run Qserv integration tests
===========================

.. code:: bash

    cd "$WORKDIR"
    git clone  https://github.com/lsst/qserv-operator
    cd qserv-operator
    ./wait-qserv-ready.sh
    ./run-integration-tests.sh