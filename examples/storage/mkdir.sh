#!/bin/bash

#  Create directory belonging to a given user
#  on all hosts

# @author Fabrice Jammes SLAC/IN2P3

set -euxo pipefail

DIR=$(cd "$(dirname "$0")"; pwd -P)
. "$DIR/env.sh"

REMOTE_DIR="/$DATA_DIR/"$NS-$INSTANCE"/qserv"
echo "Create directory $REMOTE_DIR on all nodes, if not exists"

REMOTE_USER="qserv"

for node in $MASTERS $WORKERS
do
    echo "mkdir $REMOTE_DIR on $node"
    ssh $SSH_CFG_OPT "$node" "sudo mkdir -p $REMOTE_DIR && \
    sudo chown $REMOTE_USER:$REMOTE_USER $REMOTE_DIR"
done

REMOTE_DIR="/$DATA_DIR/"$NS-$INSTANCE"/replication"
echo "Create directory $REMOTE_DIR on all Qserv master nodes, if not exists"
for node in $MASTERS
do
    echo "mkdir $REMOTE_DIR on $node"
    ssh $SSH_CFG_OPT "$node" "sudo mkdir -p $REMOTE_DIR && \
    sudo chown $REMOTE_USER:$REMOTE_USER $REMOTE_DIR"
done
