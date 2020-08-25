#!/bin/bash

# See
# https://github.com/operator-framework/operator-sdk/blob/master/doc/user-guide.md#build-and-run-the-operator

set -euxo pipefail

DIR=$(cd "$(dirname "$0")"; pwd -P)
. "$DIR/env.sh"

make install
make deploy IMG="$OP_IMAGE"
