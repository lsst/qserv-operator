#!/bin/sh

set -e
set -x

kubectl delete pod -l name=qserv-operator
kubectl delete -f /home/fjammes/go/src/github.com/lsst/qserv-operator/deploy/crds/qserv_v1alpha1_qserv_cr.yaml
kubectl apply -f /home/fjammes/go/src/github.com/lsst/qserv-operator/deploy/crds/qserv_v1alpha1_qserv_cr.yaml

