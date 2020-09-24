#!/bin/bash

# LSST Data Management System
# Copyright 2014 LSST Corporation.
# 
# This product includes software developed by the
# LSST Project (http://www.lsst.org/).
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
# 
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
# 
# You should have received a copy of the LSST License Statement and 
# the GNU General Public License along with this program.  If not, 
# see <http://www.lsstcorp.org/LegalNotices/>.

# Wait Qserv pods to be in ready state

# @author Fabrice Jammes SLAC/IN2P3

set -eux

VERBOSE=false
DIR=$(cd "$(dirname "$0")"; pwd -P)

while test $# -gt 0; do
  case "$1" in
    -v | --verbose)
      VERBOSE=true
      shift
      ;;
  esac
done

INSTANCE=$(kubectl get qservs.qserv.lsst.org -o=jsonpath='{.items[0].metadata.name}')
WORKER_COUNT=$(kubectl get qservs.qserv.lsst.org "$INSTANCE" -o=jsonpath='{.spec.worker.replicas}')

SHELL_POD="${INSTANCE}-shell"

kubectl delete pod -l "app=qserv,instance=$INSTANCE,tier=shell"
kubectl run "${INSTANCE}-shell" --image="curlimages/curl:7.70.0"  --restart=Never sleep 3600
kubectl label pod "${INSTANCE}-shell" "app=qserv" "instance=$INSTANCE" "tier=shell"
while ! kubectl wait pod --for=condition=Ready --timeout="10s" -l "app=qserv,instance=$INSTANCE"
do
  echo "Wait for Qserv pods to be ready:"
  kubectl get pod -l "app=qserv,instance=$INSTANCE"
  if [ "$VERBOSE" = true ]; then
    kubectl describe pod -l "app=qserv,instance=$INSTANCE"
  fi
done

echo "Qserv pods are ready:"
kubectl get all -l "app=qserv,instance=$INSTANCE"

for (( i=0; i<${WORKER_COUNT}; i++ ))
do
  kubectl cp "$DIR/wait-wmgr.sh" "$SHELL_POD":/tmp
  kubectl exec "$SHELL_POD" -it /tmp/wait-wmgr.sh "${INSTANCE}-worker-${i}.${INSTANCE}-worker"
done
echo "wmgr service is ready in all Qserv pods"

kubectl delete pod -l "app=qserv,instance=$INSTANCE,tier=shell"
