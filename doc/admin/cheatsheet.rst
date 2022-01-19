Kubernetes cheat sheet for Qserv
################################

Check first the `official k8s cheat sheet`_ for examples of ``kubectl`` basic commands.

Prerequisites
=============

Get access to a Kubernetes cluster running Qserv, see :ref:`installation-label`.

Interact with running pods
==========================

.. code:: shell

       # Check a qserv instance is up and running
       $ kubectl get qservs.qserv.lsst.org
       NAME        AGE
       qserv-dev   12h

       # Get instance name
       $ INSTANCE=$(kubectl get qservs.qserv.lsst.org -o=jsonpath='{.items[0].metadata.name}')
       $ echo $INSTANCE
       qserv-dev

       # Describe a pod
       $ kubectl describe pods "$INSTANCE"-worker-0

       # Get the containers list for a given pod
       $ kubectl get pods "$INSTANCE"-worker-0 -o jsonpath='{.spec.containers[*].name}'

       # Open a shell on mariadb container on worker qserv-0
       $ kubectl exec -it "$INSTANCE"-worker-0 -c mariadb bash

Access to log files
===================

.. code:: shell


       CZAR_POD="$INSTANCE-czar-0"

       # Get the xrootd container logs on pod czar-0
       kubectl logs "$CZAR_POD" -c xrootd

       # Get the previous xrootd container logs on pod czar-0
       kubectl logs "$CZAR_POD" -c xrootd -p

`Stern`_ provides advanced logging management features.

Update Qserv configuration
==========================

Update Qserv configuration by updating its related k8s configmaps.

.. code:: shell


       # List configmaps
       $ kubectl get configmaps -l app=qserv,instance="$INSTANCE"

       # Edit configmap online, i.e. directly inside etcd, the k8s database
       $ kubectl edit configmaps qserv-dev-repl-ctl-etc

       # Restart the pod using the configmap, in order to take in account the new configuration
       $ kubectl delete po "$INSTANCE"-repl-ctl-0

       # A Qserv re-install will reset the configmap, so eventually, backup the configmap locally
       $ kubectl get cm "$INSTANCE"-repl-ctl-etc -o yaml > "$INSTANCE".bck.yaml

Launch commands directly on Qserv nodes
=======================================

Check if pod worker-0 can connect to replication database and dump it configuration
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. code:: sh

    $ kubectl exec -it "$INSTANCE"-worker-0 -c repl-wrk -- mysql -h "$INSTANCE"-repl-db -u qsreplica -e "SELECT * FROM qservReplica.config;"
    +------------+---------------------------------+-----------------------------+
    | category   | param                           | value                       |
    +------------+---------------------------------+-----------------------------+
    | common     | request_buf_size_bytes          | 131072                      |
    | common     | request_retry_interval_sec      | 5                           |
    | controller | empty_chunks_dir                | /qserv/data/qserv           |
    | controller | http_server_port                | 8080                        |
    | controller | http_server_threads             | 16                          |
    | controller | job_heartbeat_sec               | 0                           |
    | controller | job_timeout_sec                 | 57600                       |
    ...

.. _Official k8s cheat sheet: https://kubernetes.io/docs/reference/kubectl/cheatsheet
.. _README: ../install
.. _Stern: https://github.com/wercker/stern

Delete a qserv instance and related storage
===========================================

.. code:: sh

    # Delete all qserv instances in current namespace
    kubectl delete qservs.qserv.lsst.org --all
    # Delete all qserv persistent volume claims in current namespace
    kubectl delete pvc -l app.kubernetes.io/managed-by=qserv-operator
