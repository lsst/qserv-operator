#!/bin/sh

set -e
set -x

DIR=$(cd "$(dirname "$0")"; pwd -P)

NS="default"

kubectl apply -f "$DIR"/deploy/crds/qserv_v1alpha1_qserv_crd.yaml --namespace="$NS"
kubectl apply -f "$DIR"/deploy/service_account.yaml --namespace="$NS"
kubectl apply -f "$DIR"/deploy/role.yaml --namespace="$NS"
kubectl apply -f "$DIR"/deploy/role_binding.yaml --namespace="$NS"
kubectl apply -f "$DIR"/deploy/operator.yaml --namespace="$NS"

echo "Run:"
echo "kubectl apply -f $DIR/deploy/crds/qserv_v1alpha1_qserv_cr.yaml --namespace='$NS'"
kubectl apply -f $DIR/deploy/crds/qserv_v1alpha1_qserv_cr.yaml --namespace="$NS"
