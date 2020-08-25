#!/usr/bin/env bash

# See
# https://github.com/operator-framework/operator-sdk/blob/master/doc/user-guide.md#build-and-run-the-operator

set -euxo pipefail

DIR=$(cd "$(dirname "$0")"; pwd -P)
. "$DIR/env.sh"

export GO111MODULE="on"
make manifests
make docker-build IMG="$OP_IMAGE"
docker push "$OP_IMAGE" || echo "WARN: unable to push qserv-operator image to Docker hub"
