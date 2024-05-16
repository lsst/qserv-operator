##########################
Fine-tune Qserv deployment
##########################

Prerequisites
=============

For all setups
--------------

-  Access to a Kubernetes v1.19+ cluster via a valid ``KUBECONFIG`` file.
-  Dynamic volume provisionning need to be available on the Kubernetes cluster (for example `kind <https://kind.sigs.k8s.io/>`__ for or
   GKE).
- `cert-manager <https://cert-manager.io/docs/installation/>`_

For a development workstation
-----------------------------

-  Ubuntu LTS is recommended
-  8 cores, 16 GB RAM, 30GB for the partition hosting docker entities
   (images, volumes, containers, etc). Use ``df`` command as below to
   find its size.

   .. code:: bash

       sudo df –sh /var/lib/docker # or /var/snap/docker/common/var-lib-docker/

-  Internet access without proxy
-  ``sudo`` access
-  Install dependencies below:

   .. code:: bash

       sudo apt-get install curl docker.io git vim

-  Add current user to docker group and restart gnome session

   .. code:: bash

       sudo usermod -a -G docker <USER>

-  Install Kubernetes locally using this `simple k8s install script <https://github.com/k8s-school/kind-helper>`__, based on
   `kind <https://kind.sigs.k8s.io/>`__.


Deploy qserv-operator
=====================

`qserv-operator` can be deployed at cluster-scope, and will then manage all Qserv instances of the k8s cluster.
It can also be deployed at namespace-scope, and will then manage Qserv instances based in the same namespace.
Installing in the same k8s cluster, a qserv-operator at cluster-scope and an other at namespace-scope is not supported.

.. note::

   If target k8s cluster does not support dynamic storage provisionning then PersistentVolumes and PersistentVolumeClaims
   must be manually created before installing Qserv. This project `Skateful <https://github.com/k8s-school/skateful` explains how to create them.

At cluster-scope
----------------

The operator will manage all Qserv instances across the cluster.

.. code:: sh

    # Deploy qserv-operator at cluster-scope in "qserv-operator-system" namespace
    RELEASE="2024.5.1-rc3"
    kubectl apply -f https://raw.githubusercontent.com/lsst/qserv-operator/$RELEASE/manifests/operator.yaml

At namespace-scope
------------------

The operator will only manage Qserv instances in the namespace where it is installed.
This setup allows running multiple instances of `qserv-operator` across a same Kubernetes cluster.
However, Qserv CustomResourceDefinitions (CRDs) are cluster-scoped, so a conflict might occurs between two operator versions if their CRDs are not the same.

.. warning::

   This setup is not compatible with a cluster-scope qserv-operator running in the same cluster.

.. code:: sh

    # Deploy qserv-operator at namespace-scope in "qserv-dev" namespace
    RELEASE="2024.5.1-rc3"
    NAMESPACE="qserv-dev"
    curl https://raw.githubusercontent.com/lsst/qserv-operator/$RELEASE/manifests/operator-ns-scoped.yaml | sed "s/<NAMESPACE>/$NAMESPACE/" | kubectl apply -f -


Deploy a qserv instance
=======================

with default settings
---------------------

Default settings below are recommended for development purpose, or continuous integration. 

.. code:: sh

    # Install a qserv instance with default settings inside a given namespace
    kubectl apply -k https://github.com/lsst/qserv-operator/manifests/base?ref=$RELEASE --namespace='<NAMESPACE>'

    # For example, at in2p3, use urls:
    # - https://github.com/lsst/qserv-operator/manifests/in2p3?ref=$RELEASE
    # - https://github.com/lsst/qserv-operator/manifests/in2p3-dev?ref=$RELEASE

with custom settings
--------------------

For production setup, Qserv install customization is handled with
`Kustomize <https://github.com/kubernetes-sigs/kustomize>`__, which is a
template engine allowing to customize kubernetes Yaml files.
``Kustomize`` is integrated with ``kubectl`` (``-k`` option).

This setup is recommended for production platforms.

Example are available, see below:

.. code:: sh

    # Install a qserv instance with custom settings
    kubectl apply -k https://github.com/lsst/qserv-operator/manifests/in2p3?ref=$RELEASE --namespace='qserv-prod'

In order to create a customized Qserv instance, create a ``Kustomize``
overlay using instructions below:

.. code:: sh

    RELEASE="2024.5.1-rc3"
    git clone --depth 1 --single-branch -b "$RELEASE" https://github.com/lsst/qserv-operator
    cd qserv-operator
    cp -r manifests/base/ manifests/<customized-overlay>

Then add custom setting, by editing ``manifests/<customized-overlay>/qserv.yaml``:

And finally create customized Qserv instance:

.. code:: sh

    kubectl apply -k manifests/<customized-overlay>/ --namespace='<namespace>'

Run Qserv integration tests
===========================

.. code:: bash

    cd "$WORKDIR"
    RELEASE="2024.5.1-rc3"
    git clone --depth 1 --single-branch -b "$RELEASE" https://github.com/lsst/qserv-operator
    cd qserv-operator
    ./tests/tools/wait-qserv-ready.sh
    ./tests/e2e/integration.sh

Undeploy a Qserv instance
=========================

First list all Qserv instances running in a given namespace

.. code:: sh

    kubectl get qserv -n "<namespace>"

It will output something like:

::

    NAME    CZARS   INGEST-DB   REPL-CTL   REPL-DB   WORKERS   XROOTD   AGE
    qserv   1/1     1/1         1/1        1/1       2/2       2/2      2d10h


Then delete this Qserv instance

.. code:: sh

    kubectl delete qserv qserv -n "<namespace>"

To delete all Qserv instances inside a namespace:

.. code:: sh

    kubectl delete qserv --all -n "<namespace>"

Qserv storage will remain untouch by this operation, and a restarting Qserv instance will use the existing storage. To remove the existing Qserv storage and re-initialize Qserv databases from scratch, run:

.. code:: sh

    # Delete all qserv persistent volume claims in current namespace
    kubectl delete pvc -l app.kubernetes.io/managed-by=qserv-operator
