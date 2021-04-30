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
while true; do
  if kubectl wait qserv --timeout=10s --for=condition=available "$INSTANCE"
  then
        kubectl get qservs.qserv.lsst.org -o jsonpath='{.items[*].status}'
        break
    else
        echo "Wait for Qserv to start"
        kubectl get pods,pvc -l app=qserv
        kubectl get pv
        # See https://stackoverflow.com/questions/67122591/display-logs-of-an-initcontainer-running-inside-github-actions
        # INITDB_WAITING=$(kubectl get pods qserv-worker-0 -o jsonpath='{$.status.initContainerStatuses[0].state.waiting}')
        # INITDB_WAITING=$(kubectl get pods qserv-worker-0 -o jsonpath='{$.status.phase}')
        echo "initdb state"
        POD_WORKER=$(kubectl get pods -l statefulset.kubernetes.io/pod-name=qserv-worker-0 -o jsonpath='{.items[0].metadata.name}')
        if [ -n "$POD_WORKER" ]; then
          INITDB_RUNNING=$(kubectl get pods $POD_WORKER -o jsonpath='{$.status.initContainerStatuses[0].state.running}')
          if [ -n "$INITDB_RUNNING" ]; then
            echo "initdb logs for $POD_WORKER"
            kubectl logs qserv-worker-0 -c initdb
          fi
        fi
    fi
done
echo "Qserv pods are ready"
kubectl get pods -l app=qserv
kubectl describe statefulsets.apps qserv-czar
kubectl get pods -A
kubectl describe pod qserv-czar-0
kubectl get qservs.qserv.lsst.org -o jsonpath='{.items[*].status}'
