#!/bin/sh

# See
# https://github.com/operator-framework/operator-sdk/blob/master/doc/user-guide.md#build-and-run-the-operator

set -eux

DIR=$(cd "$(dirname "$0")"; pwd -P)
. "$DIR/env.sh"

GO111MODULE="on" operator-sdk build "$OP_IMAGE"
sed "s|REPLACE_IMAGE|$OP_IMAGE|g" "$DIR/deploy/operator.yaml.tpl" \
    > "$DIR/deploy/operator.yaml"
docker push "$OP_IMAGE" || echo "WARN: unable to push image to Docker hub"
