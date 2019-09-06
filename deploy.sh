#!/bin/sh

set -e
set -x

DIR=$(cd "$(dirname "$0")"; pwd -P)

kubectl apply -f "$DIR"/deploy/crds/qserv_v1alpha1_qserv_crd.yaml
kubectl apply -f "$DIR"/deploy/service_account.yaml
kubectl apply -f "$DIR"/deploy/role.yaml
kubectl apply -f "$DIR"/deploy/role_binding.yaml
kubectl apply -f "$DIR"/deploy/operator.yaml

echo "Run:"
echo "kubectl apply -f $DIR/deploy/crds/qserv_v1alpha1_qserv_cr.yaml"
kubectl apply -f $DIR/deploy/crds/qserv_v1alpha1_qserv_cr.yaml
