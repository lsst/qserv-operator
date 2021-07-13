###################
Retrieve core files
###################

Set core file path on infrastructure
====================================

This `article on Kubernetes and core file managemnt <https://medium.com/faun/handling-core-dumps-in-kubernetes-clusters-in-gcp-b1b2a54c25dc>`__
and this `other one <https://medium.com/@shuanglu1993/how-to-generate-coredump-for-containers-running-with-k8s-1a3f4a7e75b2>`__ provide some tracks about core file management.

First, on your Kubernetes nodes, set the core file path, named `<core-path>`, in example below it is `/tmp/coredump`, and make it writable to all container users which can potentially create core files:

.. code:: bash

   COREPATH=/tmp/coredump
   sudo mkdir $COREPATH
   sudo chmod 777 $COREPATH
   sudo sh -c 'echo "$COREPATH/core.%e.%p.%h.%t" > /proc/sys/kernel/core_pattern'

.. note::

   For `kind <https://kind.sigs.k8s.io>`_ users, on development workstation, above command must be run on host, it will be then propagated to k8s node and to pods.

Store core files in a persistent storage
========================================

:doc:`Install Qserv operator <../install/main>` and then install a Qserv instance dedicated to development:

.. code:: bash
   
   kubectl apply -k https://github.com/lsst/qserv-operator/manifests/dev

Core files produced by every Qserv binaries will be stored and available.

For additional information, see :ref:`fine-tune-qserv-dev-instance`

Download core files
===================

Core files will be available on the node running the pods, inside the `<core-path>` directory.

.. note::

   For `kind <https://kind.sigs.k8s.io>`_ users, use `docker` command to get core file on the workstation (see demo below). For bare-metal Kubernetes clusters, `scp` or `rsync` should work fine.

Demo
====

This demo rely on a Kubernetes cluster based on kind and the `qserv-operator`:

.. code:: bash

   # Create a directory to store core files on the k8s node (kind-specific)
   docker exec -it -- kind-control-plane sh -c "mkdir -p /tmp/coredump && chmod 777 /tmp/coredump"

   # Install Qserv
   kubectl apply -k qserv-operator/manifests/dev

   # Check Qserv is running
   kubectl get pods -o wide
   NAME                              READY   STATUS    RESTARTS   AGE     IP            NODE              
   qserv-dev-czar-0                  3/3     Running   0          11s     10.244.0.54   kind-control-plane
   qserv-dev-repl-ctl-0              1/1     Running   0          11s     10.244.0.46   kind-control-plane
   qserv-dev-repl-db-0               1/1     Running   0          11s     10.244.0.56   kind-control-plane
   qserv-dev-worker-0                5/5     Running   0          11s     10.244.0.53   kind-control-plane
   qserv-dev-worker-1                5/5     Running   0          11s     10.244.0.51   kind-control-plane
   qserv-dev-worker-2                5/5     Running   0          11s     10.244.0.55   kind-control-plane
   qserv-dev-xrootd-redirector-0     2/2     Running   0          11s     10.244.0.44   kind-control-plane
   qserv-dev-xrootd-redirector-1     2/2     Running   0          10s     10.244.0.45   kind-control-plane
   qserv-dev-xrootd-redirector-2     2/2     Running   0          10s     10.244.0.47   kind-control-plane
   qserv-dev-xrootd-redirector-3     2/2     Running   0          10s     10.244.0.48   kind-control-plane
   qserv-operator-5467b89db4-hbwgc   1/1     Running   0          149m    10.244.0.5    kind-control-plane

   # Kill replication controller
   kubectl exec -it qserv-dev-repl-ctl-0 -- bash
   bash-4.2$ ps -ef
   UID        PID  PPID  C STIME TTY          TIME CMD
   qserv        1     0  0 11:30 ?        00:00:00 /bin/sh /config-start/start.sh
   qserv        9     1  0 11:30 ?        00:00:00 qserv-replica-master-http --worker-evict-timeout=3600 --health-probe-interval=120 --replication-interval=1200 --config=mysql://qsreplica:@qserv-dev-repl-db:3306/qservReplica --qserv-db-password=CHANGEME
   qserv      100     0  0 11:38 pts/0    00:00:00 bash
   qserv      112   100  0 11:38 pts/0    00:00:00 ps -ef
   bash-4.2$ kill -s SIGSEGV  9
   bash-4.2$ command terminated with exit code 137

   # List and retrieve core file (kind-specific)
   docker ls  docker exec -it kind-control-plane ls /tmp/coredump
   core.qserv-replica-m.9.qserv-dev-repl-ctl-0.1597318703

   # Retrieve corefile locally (docker cp does not work because /tmp is managed by tmpfs in kind)
   docker exec kind-control-plane tar Ccf "/tmp/coredump" - . | tar Cxf . -
   ls 
   core.qserv-replica-m.9.qserv-dev-repl-ctl-0.1597318703

#################################################
Debug manually a process inside a Qserv container
#################################################

Install a Qserv instance dedicated to development
=================================================

:doc:`Install Qserv operator <../install/main>` and then install a Qserv instance dedicated to development:

.. code:: bash
   
   kubectl apply -k https://github.com/lsst/qserv-operator/manifests/dev

Demo
====

.. code:: bash

    kubectl exec -it qserv-dev-repl-ctl-0 bash

    bash-4.2$ gdb /qserv/bin/qserv-replica-master-http
    GNU gdb (GDB) Red Hat Enterprise Linux 7.12.1-48.el7
    ...
    Reading symbols from /qserv/bin/qserv-replica-master-http...done.

    (gdb) run --config=mysql://qsreplica@lsst-qserv-master01:23306/qservReplica --instance-id=qserv-prod --qserv-db-password=xxx --auth-key=xxx --debug

.. _fine-tune-qserv-dev-instance:

######################################
Fine-tune a Qserv development instance
######################################

Pre-requisites
==============

First, download `qserv-operator` locally

.. code:: bash
   
   git clone https://github.com/lsst/qserv-operator

Core path
=========

It is possible to set the core path easily by editing the `corepath` parameter in file `qserv-operator/manifests/dev/qserv.yaml`

.. code:: yaml

   apiVersion: qserv.lsst.org/v1alpha1
   kind: Qserv
   metadata:
   name: qserv
   spec:
   devel:
      corepath: "<core-path>"


Manual debugging with gdb
=========================

It is possible to set the component(s) to debug by editing the `debug` parameters in file `qserv-operator/manifests/dev/qserv.yaml`

.. code:: yaml

   apiVersion: qserv.lsst.org/v1alpha1
   kind: Qserv
   metadata:
   name: qserv
   spec:
   ...
   replication:
     debug: "repl-ctl"

Values for the debug parameter are:
* `repl-ctl`: replication controller start in debug mode.
* `repl-wrk`: all replication worker start in debug mode.
* `all`: both replication controller and replication workers start in debug mode.

In above example, replication controller will not start, so that user can open an interactive shell inside the container,
start the replication controller process in debug mode and perform debugging operation. The container won't restart if the replication controller crashes.


Re-install Qserv
================

Once file `qserv-operator/manifests/dev/qserv.yaml` is ready, (re-)install Qserv in the current namespace

.. code:: bash
   
   kubectl apply -k qserv-operator/manifests/dev
