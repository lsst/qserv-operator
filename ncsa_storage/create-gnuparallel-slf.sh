#!/bin/bash

set -e
set -x

# Create a *.slf file used by GNU parallel
# to access Qserv cluster

# @author  Fabrice Jammes, IN2P3

DIR=$(cd "$(dirname "$0")"; pwd -P)

. "$DIR/env.sh"

if [ ! -e "$SSH_CFG" ]
then
    echo "WARN: non-existing $SSH_CFG"
fi

rm -f "$PARALLEL_SSH_CFG"
for n in $MASTER $WORKERS
do
    echo "ssh $SSH_CFG_OPT $n" >> "$PARALLEL_SSH_CFG"
done 
