#!/bin/sh

set -e
set -x

DIR=$(cd "$(dirname "$0")"; pwd -P)

kubectl delete -f "$DIR"/deploy/crds
kubectl delete -f "$DIR"/deploy
