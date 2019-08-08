#!/bin/sh

set -e
set -x

DIR=$(cd "$(dirname "$0")"; pwd -P)
. "$DIR"/../env.sh

docker pull "$QSERV_IMAGE" 

kind  load docker-image "$OP_IMAGE" 
kind load docker-image "$QSERV_IMAGE"
