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

set -e

DIR=$(cd "$(dirname "$0")"; pwd -P)

. "$DIR/env.sh"

SHELL_POD="${INSTANCE}-shell"

echo "Wait for Qserv pods to be ready"
kubectl run "${INSTANCE}-shell" --image=alpine  --restart=Never sleep 3600
kubectl label pod "${INSTANCE}-shell" "app=qserv" "instance=$INSTANCE" "tier=shell"
kubectl wait pod --for=condition=Ready --timeout="-1s" -l "app=qserv,instance=$INSTANCE"
kubectl cp "$DIR/wait-wmgr.sh" "$SHELL_POD":/root
kubectl exec "$SHELL_POD" -it /root/wait-wmgr.sh example-qserv-worker-0.example-qserv-worker
kubectl exec "$SHELL_POD" -it /root/wait-wmgr.sh example-qserv-worker-1.example-qserv-worker
kubectl exec "$SHELL_POD" -it /root/wait-wmgr.sh example-qserv-worker-2.example-qserv-worker
kubectl delete pod -l "app=qserv,instance=$INSTANCE,tier=shell"
