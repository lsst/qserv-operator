#!/bin/bash

#  Create directory belonging to a given user
#  on all hosts

# @author Fabrice Jammes SLAC/IN2P3

set -e
set -x

DIR=$(cd "$(dirname "$0")"; pwd -P)
. "$DIR/env.sh"

REMOTE_DIR="/qserv/"$INSTANCE"/qserv"
echo "Create directory $REMOTE_DIR on all nodes, if not exists"

REMOTE_USER="qserv"

for node in $MASTER $WORKERS
do
    echo "mkdir $REMOTE_DIR on $node"
    ssh $SSH_CFG_OPT "$node" "sudo mkdir -p $REMOTE_DIR && \
    sudo chown $REMOTE_USER:$REMOTE_USER $REMOTE_DIR"
done
