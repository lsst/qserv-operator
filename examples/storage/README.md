# Setup local storage for Qserv PersistentVolumes at NCSA

## Pre-requisites

Clone qserv operator repository and go to storage management directory

```shell
git clone https://github.com/lsst/qserv-operator.git
cd qserv-operator/ncsa_storage
git checkout tickets/DM-21824
```

## 1 - Create data directories

Create `/qserv/qserv-dev/qserv` on each nodes, including master.

Create `/qserv/qserv-dev/repl` on master node only.

## 2 - Create StorageClass, PersistentVolumes and PersistentVolumesClaims

```shell
kubectl apply -n qserv-dev -f out/
```
