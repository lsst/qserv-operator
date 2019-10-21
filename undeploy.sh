#!/bin/bash

set -euxo pipefail

DIR=$(cd "$(dirname "$0")"; pwd -P)

kubectl delete all,configmaps -l app=qserv
kubectl delete qserv --all

PVCS=$(kubectl get pvc -l app=qserv -o go-template='{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}')

for pvc in $PVCS
do
    ANNOTS=$(kubectl get pvc "$pvc" -o go-template='{{.metadata.annotations}}')
    STRING="volume.beta.kubernetes.io/storage-provisioner:kubernetes.io/host-path"
    case "$ANNOTS" in
    *"$STRING"* ) kubectl delete pvc "$pvc"
    esac
done

kubectl delete -f "$DIR"/deploy/crds/qserv_v1alpha1_qserv_crd.yaml
kubectl delete -f "$DIR"/deploy