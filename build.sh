#!/bin/sh

# See
# https://github.com/operator-framework/operator-sdk/blob/master/doc/user-guide.md#build-and-run-the-operator

set -e
set -x

DIR=$(cd "$(dirname "$0")"; pwd -P)
. "$DIR/env.sh"

GO111MODULE="on" operator-sdk build "$OP_IMAGE"
sed "s|REPLACE_IMAGE|$OP_IMAGE|g" "$DIR/deploy/operator.yaml.tpl" \
    > "$DIR/deploy/operator.yaml"
