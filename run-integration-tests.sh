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
# Launch Qserv multinode tests on Swarm cluster

# @author Fabrice Jammes SLAC/IN2P3

set -eux

INSTANCE=$(kubectl get qservs.qserv.lsst.org -o=jsonpath='{.items[0].metadata.name}')
WORKER_COUNT=$(kubectl get qservs.qserv.lsst.org "$INSTANCE" -o=jsonpath='{.spec.worker.replicas}')
CSS_INFO=""

# Build CSS input data
upper_id=$((WORKER_COUNT-1))
for i in $(seq 0 "$upper_id");
do
    CSS_INFO="${CSS_INFO}CREATE NODE worker${i} type=worker port=5012 \
    host=${INSTANCE}-worker-${i}.${INSTANCE}-worker; "
done

kubectl exec "${INSTANCE}-czar-0" -c wmgr -- su qserv -l -c ". /qserv/stack/loadLSST.bash && \
    setup qserv_distrib -t qserv-dev && \
    echo \"$CSS_INFO\" | qserv-admin.py -c mysql://qsmaster@127.0.0.1:3306/qservCssData && \
    qserv-test-integration.py -V DEBUG"

