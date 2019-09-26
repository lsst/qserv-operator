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

# Wait Qserv statefulset to be in running state

# @author Fabrice Jammes SLAC/IN2P3

set -e

DIR=$(cd "$(dirname "$0")"; pwd -P)

. "$DIR/env.sh"

echo "Wait for Qserv statefulsets to be in running state"


GO_TPL="{{if .status.readyReplicas}}\
    .status.readyReplicas is set \
    {{end}}"

for sf in "${INSTANCE}-czar" "${INSTANCE}-worker" "${INSTANCE}-xrootd-redirector"
do
    echo -n "Wait for statefulset '$sf' to exist"
    until kubectl get statefulset/"$sf" > /dev/null 2>&1 
    do
        sleep 2
        echo -n '.'
    done
    echo

    echo -n "Wait for statefulset '$sf' to start first pod"
    until [ -n "$READY" ]
    do
        READY=$(kubectl get statefulset "$sf" -o go-template --template "$GO_TPL")
        sleep 2
	echo -n '.'
    done
    echo

    echo -n "Wait for statefulset '$sf' to start all pods"
    GO_TPL="{{if and (eq .spec.replicas .status.replicas) \
        (eq .status.replicas .status.readyReplicas) \
        (eq .status.currentRevision .status.updateRevision)}}true{{end}}"
    until [ -n "$STARTED" ]
    do
        STARTED=$(kubectl get statefulset "$sf" -o go-template --template "$GO_TPL")
        sleep 2
	echo -n '.'
    done
    echo
    echo "Statefulset '$sf' ready"
done
