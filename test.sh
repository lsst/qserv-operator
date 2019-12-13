#!/bin/sh

# See
# https://github.com/operator-framework/operator-sdk/blob/master/doc/user-guide.md#build-and-run-the-operator

set -eux

DIR=$(cd "$(dirname "$0")"; pwd -P)

go test "$DIR/pkg/..." "$DIR/cmd/..." -cover
