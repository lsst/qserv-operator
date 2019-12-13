#!/bin/sh

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

# Creates K8s Volumes and Claims for Master and Workers

# @author Benjamin Roziere, IN2P3
# @author Fabrice Jammes, IN2P3

set -eux

DIR=$(cd "$(dirname "$0")"; pwd -P)

. "$DIR/env.sh"

usage() {
    cat << EOD

    Usage: $(basename "$0") <hostPath>

    Available options:
      -h          this message

      Generate yaml files for Qserv PersistentVolumeClaims and PersistentVolumes.
      These PersistentVolumes are based on local-storage class
      and use <hostPath> as a local storage on each nodes.

EOD
}

while getopts hp: c ; do
    case $c in
        h) usage; exit 0 ;;
        \?) usage ; exit 2 ;;
    esac
done
shift "$((OPTIND-1))"

if [ $# -ne 1  ] ; then
    usage
    exit 2
fi

STORAGE_PATH="$1"

YAML_OUT_DIR=$DIR/out
rm -rf $YAML_OUT_DIR
mkdir $YAML_OUT_DIR

cp "$DIR/manifests/storageclass-qserv.yaml" "$YAML_OUT_DIR"

PVC_PREFIX="qserv-data-${INSTANCE}"

DATA_PATH="$STORAGE_PATH/${INSTANCE}/data"

echo "Creating persistent volumes and claims for Qserv czars"
COUNT=0
for host in $MASTERS;
do
    OPT_HOST="-H $host"
    PVC_NAME="${PVC_PREFIX}-czar-${COUNT}"
    "$DIR"/yaml-builder.py -p "$DATA_PATH" -n "$PVC_NAME" $OPT_HOST -o "$YAML_OUT_DIR" -i "$INSTANCE"
    COUNT=$((COUNT+1))
done

echo "Creating persistent volumes and claims for Qserv"
COUNT=0
for host in $WORKERS;
do
    OPT_HOST="-H $host"
    PVC_NAME="${PVC_PREFIX}-worker-${COUNT}"
    "$DIR"/yaml-builder.py -p "$DATA_PATH" -n "$PVC_NAME" $OPT_HOST -o "$YAML_OUT_DIR" -i "$INSTANCE"
    COUNT=$((COUNT+1))
done

echo "Creating persistent volumes and claims for Replication Database"
OPT_HOST="-H $REPL_DB_HOST"
PVC_NAME="${PVC_PREFIX}-repl-db-0"
DATA_PATH="$STORAGE_PATH/${INSTANCE}/replication"
"$DIR"/yaml-builder.py -p "$DATA_PATH" -n "$PVC_NAME" $OPT_HOST -o "$YAML_OUT_DIR" -i "$INSTANCE"
