# Kubernetes cheat sheet for Qserv

[Official k8s cheat sheet](https://kubernetes.io/docs/reference/kubectl/cheatsheet)

## Pre-requisites

Get access to a Kubernetes cluster running Qserv, see [README](../README.md)

## Interact with running pods

```shell

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
```

## Access to log files

```shell

    CZAR_POD="$INSTANCE-czar-0"

    # Get the xrootd container logs on pod czar-0
    kubectl logs "$CZAR_POD" -c xrootd

    # Get the previous xrootd container logs on pod czar-0
    kubectl logs "$CZAR_POD" -c xrootd -p
```

[Stern](https://github.com/wercker/stern) provides advanced logging management features.

## Update Qserv configuration

Update Qserv configuration by updating its related k8s configmaps.

```shell

    # List configmaps
    $ kubectl get configmaps -l app=qserv,instance="$INSTANCE"

    # Edit configmap online, i.e. directly inside etcd, the k8s database
    $ kubectl edit configmaps qserv-dev-repl-ctl-etc

    # Restart the pod using the configmap, in order to take in account the new configuration
    $ kubectl delete po "$INSTANCE"-repl-ctl-0

    # A Qserv re-install will reset the configmap, so eventually, backup the configmap locally
    $ kubectl get cm "$INSTANCE"-repl-ctl-etc -o yaml > "$INSTANCE".bck.yaml
```

## Launch commands directly on Qserv nodes

### Check if pod worker-0 can connect to replication database and dump it configuration

```shell

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
    | controller | num_threads                     | 64                          |
    | controller | request_timeout_sec             | 57600                       |
    | database   | qserv_master_host               | qserv-dev-czar              |
    | database   | qserv_master_name               | qservMeta                   |
    | database   | qserv_master_port               | 3306                        |
    | database   | qserv_master_services_pool_size | 4                           |
    | database   | qserv_master_tmp_dir            | /qserv/data/injest          |
    | database   | qserv_master_user               | qsmaster                    |
    | database   | services_pool_size              | 32                          |
    | worker     | data_dir                        | /qserv/data/mysql           |
    | worker     | db_port                         | 3306                        |
    | worker     | db_user                         | root                        |
    | worker     | fs_buf_size_bytes               | 4194304                     |
    | worker     | fs_port                         | 25001                       |
    | worker     | loader_port                     | 25002                       |
    | worker     | loader_tmp_dir                  | /qserv/data/ingest          |
    | worker     | num_fs_processing_threads       | 32                          |
    | worker     | num_loader_processing_threads   | 16                          |
    | worker     | num_svc_processing_threads      | 16                          |
    | worker     | svc_port                        | 25000                       |
    | worker     | technology                      | FS                          |
    | xrootd     | auto_notify                     | 1                           |
    | xrootd     | host                            | qserv-dev-xrootd-redirector |
    | xrootd     | port                            | 1094                        |
    | xrootd     | request_timeout_sec             | 600                         |
    +------------+---------------------------------+-----------------------------+
```

### Launch a SQL query against local mariadb on worker-0

```shell
    $ kubectl exec -it "$INSTANCE"-worker-0 -c mariadb bash
    . /qserv/stack/loadLSST.bash
    setup qserv_distrib -t qserv-dev
    mysql --socket /qserv/data/mysql/mysql.sock --user=root --password="CHANGEME" -e "SHOW PROCESSLIST"
    exit
```

### Install debug tools inside a Qserv pod

```shell
    # NOTE: tools will be removed when pod restart
    # WARN: being root inside a container is insecure but is useful for development mode
    $ kubectl exec -it "$INSTANCE"-worker-0 -c mariadb bash
    # Eventually define proxy if needed
    [root@qserv-dev-worker-0 qserv]$ export https_proxy="http://ccqservproxy.in2p3.fr:3128"
    [root@qserv-dev-worker-0 qserv]$ yum install gdb bind-utils
    exit
```

## Run a pod inside the Qserv network

```shell

    # Start a pod and install debugging tools inside it
    $ kubectl run shell  --image=ubuntu  --restart=Never -- sh -c "apt-get -y update && apt-get install -y curl && sleep 3600"

    # Open a shell inside the pod
    $ kubectl exec -it shell
    root@shell:/$ curl http://qserv-dev-repl-ctl:8080
    <html>
    <head><title>404 Not Found</title></head>
    <body style="background-color:#E6E6FA">
    <h1>404 Not Found</h1>
    </body>
    </html>

    # Delete the pod
    kubectl delete pod shell
```

## Run a Qserv client pod and launch SQL queries from it

```shell

    # Start a mysql pod
    kubectl run qservclient --image=mariadb  --restart=Never -- sleep 3600

    # Open a shell and launch a query from the pod
    kubectl exec -it qservclient bash
    root@qservclient:/$ mysql --host qserv-prod-czar --port 4040 --user qsmaster gaia_dr2_00 -e "SHOW TABLES;"
    root@qservclient:/$ exit

    # Launch a query from the pod
    kubectl exec -it qservclient -- mysql --host qserv-prod-czar --port 4040 --user qsmaster -e "SELECT source_id FROM gaia_dr2_00.gaia_source WHERE source_id=4295806720;"
    +------------+
    | source_id  |
    +------------+
    | 4295806720 |
    +------------+
```

## Interact with storage

```shell
    # Get persistent volume claims for Qserv pods
    kubectl get pvc -l app=qserv

    # Get persistent volumes for Qserv pods
    kubectl get pv -l app=qserv
```