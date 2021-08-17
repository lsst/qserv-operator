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

Open a shell in the `debugger` container and list Pod's full processes list:

.. code:: bash

   kubectl exec -it qserv-worker-0 -c debugger -- bash
   [root@qserv-worker-0 /]# ps -ef
   UID          PID    PPID  C STIME TTY          TIME CMD
   65535          1       0  0 13:06 ?        00:00:00 /pause
   root          20       0  0 13:06 ?        00:00:00 /bin/sh /config-start/start.sh
   1000          28      20  0 13:06 ?        00:00:02 mysqld
   1000          60       0  0 13:06 ?        00:00:00 sleep infinity
   root          67       0  0 13:06 ?        00:00:00 /bin/sh /config-start/start.sh
   root          74      67  0 13:06 ?        00:00:00 su qserv -c sh /config-start/wmgr.sh
   1000          75      74  0 13:06 ?        00:00:00 sh /config-start/wmgr.sh
   1000          82      75  0 13:06 ?        00:00:00 python /qserv/stack/stack/miniconda3-py37_4.8.2-1eb92eb/Linux64/qserv/2021.7.1-rc1+2c8521dd9c/bin/qservWmg
   root          84       0  0 13:06 ?        00:00:00 /bin/sh /config-start/start.sh -S cmsd
   root          92      84  0 13:06 ?        00:00:00 su qserv -c /config-start/xrd.sh -S cmsd
   1000          93      92  0 13:06 ?        00:00:00 /bin/sh /config-start/xrd.sh -S cmsd
   1000          99      93  0 13:06 ?        00:00:00 cmsd -c /config-etc/xrootd.cf -n worker -I v4 -l @libXrdSsiLog.so -+xrdssi /config-etc/xrdssi.cf
   root         221       0  0 13:06 ?        00:00:00 /bin/sh /config-start/start.sh
   root         232     221  0 13:06 ?        00:00:00 su qserv -c /config-start/xrd.sh -S xrootd
   1000         233     232  0 13:06 ?        00:00:00 /bin/sh /config-start/xrd.sh -S xrootd
   1000         238     233  0 13:06 ?        00:00:00 xrootd -c /config-etc/xrootd.cf -n worker -I v4 -l @libXrdSsiLog.so -+xrdssi /config-etc/xrdssi.cf
   root         403       0  0 13:06 pts/0    00:00:00 /usr/bin/bash
   root         689       0  0 13:19 pts/1    00:00:00 bash
   root         761     689  0 13:22 pts/1    00:00:00 ps -ef

Attach `gdb` to `xrootd` process:

.. code:: bash

   # Helper to display gdb command line
   [root@qserv-worker-0 /]# debugtools 238
   2021/08/09 13:24:37 Path to executable: /proc/238/root/qserv/stack/stack/miniconda3-py37_4.8.2-1eb92eb/Linux64/xrootd/affinity-flex-hash-g5b015dcebc/bin/xrootd
   2021/08/09 13:24:37 gdb command-line: gdb -iex "set sysroot /proc/238/root" -iex "set auto-load safe-path /proc/238/root" -p 238 /proc/238/root/qserv/stack/stack/miniconda3-py37_4.8.2-1eb92eb/Linux64/xrootd/affinity-flex-hash-g5b015dcebc/bin/xrootd
   [root@qserv-worker-0 /]# gdb -iex "set sysroot /proc/238/root" -iex "set auto-load safe-path /proc/238/root" -p 238 /proc/238/root/qserv/stack/stack/miniconda3-py37_4.8.2-1eb92eb/Linux64/xrootd/affinity-flex-hash-g5b015dcebc/bin/xrootd
   ...
   Loaded symbols for /proc/238/root/qserv/stack/conda/miniconda3-py37_4.8.2/envs/lsst-scipipe-1eb92eb/lib/./libicui18n.so.67
   0x00007f3e15fd3afb in do_futex_wait.constprop.1 () from /proc/238/root/lib64/libpthread.so.0
   (gdb) bt
   #0  0x00007f3e15fd3afb in do_futex_wait.constprop.1 () from /proc/238/root/lib64/libpthread.so.0
   #1  0x00007f3e15fd3b8f in __new_sem_wait_slow.constprop.0 () from /proc/238/root/lib64/libpthread.so.0
   #2  0x00007f3e15fd3c2b in sem_wait@@GLIBC_2.2.5 () from /proc/238/root/lib64/libpthread.so.0
   #3  0x00005636d5b98959 in Wait (this=<optimized out>)
      at /qserv/stack/stack/miniconda3-py37_4.8.2-1eb92eb/EupsBuildDir/Linux64/xrootd-affinity-flex-hash-g5b015dcebc/xrootd-affinity-flex-hash-g5b015dcebc/src/./XrdSys/XrdSysPthread.hh:421
   #4  mainAccept(void*) ()
      at /qserv/stack/stack/miniconda3-py37_4.8.2-1eb92eb/EupsBuildDir/Linux64/xrootd-affinity-flex-hash-g5b015dcebc/xrootd-affinity-flex-hash-g5b015dcebc/src/Xrd/XrdMain.cc:129
   #5  0x00005636d5b8f5e2 in main (argc=<optimized out>, argv=<optimized out>)
      at /qserv/stack/stack/miniconda3-py37_4.8.2-1eb92eb/EupsBuildDir/Linux64/xrootd-affinity-flex-hash-g5b015dcebc/xrootd-affinity-flex-hash-g5b015dcebc/src/Xrd/XrdMain.cc:213


Lots of additional debugging tools are available inside the `debugtools <https://github.com/k8s-school/debugtools>`_ image,
Check the `debugtools documentation <https://github.com/k8s-school/debugtools/blob/main/README.md>`_ for additional information.

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
* `all`: both replication controller and replication workers start in debug mode.
* `repl-ctl`: replication controller start in debug mode.
* `repl-wrk`: all replication worker start in debug mode.

In above example, replication controller will not start, so that user can open an interactive shell inside the container,
start the replication controller process in debug mode and perform debugging operation. The container won't restart if the replication controller crashes.


Re-install Qserv
================

Once file `qserv-operator/manifests/dev/qserv.yaml` is ready, (re-)install Qserv in the current namespace

.. code:: bash

   kubectl apply -k qserv-operator/manifests/dev
