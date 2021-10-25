#!/bin/bash

# Deploy Qserv operator:
#   must be used from inside qserv-operator git repository, install local version of qserv-operator

set -euxo pipefail

DIR=$(cd "$(dirname "$0")"; pwd -P)

. "$DIR/env.build.sh"
make deploy
$DIR/tests/tools/wait-operator-ready.sh
