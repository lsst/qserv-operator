#!/bin/sh

set -e
set -x

DIR=$(cd "$(dirname "$0")"; pwd -P)

kubectl delete all,configmaps,pv,pvc -l app=qserv
kubectl delete qserv --all

kubectl delete -f "$DIR"/deploy/crds/qserv_v1alpha1_qserv_crd.yaml
kubectl delete -f "$DIR"/deploy

