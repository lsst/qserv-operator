#!/bin/sh

set -e
set -x

DIR=$(cd "$(dirname "$0")"; pwd -P)

IMAGES="qserv/qserv:4d93c1e qserv/replica:tools-w.2018.16-1126-g1849b21 \
        qserv/qserv-operator:v0.0.3 mariadb:10.2.16"

for img in $IMAGES
do
  kind load docker-image "$img" 
done
