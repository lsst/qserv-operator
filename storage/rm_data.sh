#!/bin/bash

#  Create directory belonging to a given user
#  on all hosts

# @author Fabrice Jammes SLAC/IN2P3

#set -e
set -uxo pipefail

DIR=$(cd "$(dirname "$0")"; pwd -P)
. "$DIR/env.sh"

TARGET_DIR="${DATA_DIR}/${NS}-${INSTANCE}/replication"
parallel --tag --nonall --slf "$PARALLEL_SSH_MASTER" "echo '$PASSWORD' | sudo -S su qserv sh -c 'rm -fr ${TARGET_DIR} && mkdir ${TARGET_DIR} && chown qserv:qserv ${TARGET_DIR}'"
TARGET_DIR="${DATA_DIR}/${NS}-${INSTANCE}/ingest"
parallel --tag --nonall --slf "$PARALLEL_SSH_MASTER" "echo '$PASSWORD' | sudo -S su qserv sh -c 'rm -fr ${TARGET_DIR} && mkdir ${TARGET_DIR} && chown qserv:qserv ${TARGET_DIR}'"
TARGET_DIR="${DATA_DIR}/${NS}-${INSTANCE}/qserv"
parallel --tag --nonall --slf "$PARALLEL_SSH_CFG" "echo '$PASSWORD' | sudo -S su qserv sh -c 'rm -fr ${TARGET_DIR} && mkdir ${TARGET_DIR} && chown qserv:qserv ${TARGET_DIR}'"
