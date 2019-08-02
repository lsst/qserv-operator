#!/bin/sh

set -e
set -x

DIR=$(cd "$(dirname "$0")"; pwd -P)

kubectl delete pod -l name=qserv-operator
kubectl delete -f "$DIR"/deploy/crds/qserv_v1alpha1_qserv_cr.yaml
kubectl apply -f "$DIR"/deploy/crds/qserv_v1alpha1_qserv_cr.yaml

