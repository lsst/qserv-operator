#!/bin/bash

set -euxo pipefail

DIR=$(cd "$(dirname "$0")"; pwd -P)
. "$DIR/env.sh"

kubectl apply -f "$DIR"/deploy/crds/qserv_v1alpha1_qserv_crd.yaml --namespace="$NS"
kubectl apply -f "$DIR"/deploy/service_account.yaml --namespace="$NS"
kubectl apply -f "$DIR"/deploy/role.yaml --namespace="$NS"
kubectl apply -f "$DIR"/deploy/role_binding.yaml --namespace="$NS"
kubectl apply -f "$DIR"/deploy/operator.yaml --namespace="$NS"

echo "----------------------------------"
echo "Run command below to deploy Qserv:"
echo "----------------------------------"
echo "kubectl apply -k $DIR/base --namespace='$NS'"
