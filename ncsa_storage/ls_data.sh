#!/bin/bash

#  Create directory belonging to a given user
#  on all hosts

# @author Fabrice Jammes SLAC/IN2P3

set -e
set -x

DIR=$(cd "$(dirname "$0")"; pwd -P)
. "$DIR/env.sh"

parallel --tag --nonall --slf "$PARALLEL_SSH_CFG" "ls -Rtl /qserv/${INSTANCE}"