#!/bin/bash

# Deploy Qserv operator:
#   - when used from inside qserv-operator git repository, install local version of qserv-operator
#   - when used in standalone mode , clone qserv-operator to temporary directory and launch inner deploy.sh

set -euxo pipefail

DIR=$(cd "$(dirname "$0")"; pwd -P)

OPERATOR='qserv-operator'
GIT_REF='tickets/DM-26295'
REPO_NAME=''
if cd $DIR && git rev-parse --is-inside-work-tree 2>/dev/null; then
  REMOTE_REPO_URL=$(git --git-dir=$DIR/.git remote get-url origin)
  REPO_NAME=$(basename -s .git $REMOTE_REPO_URL)
fi

GIT_REF="tickets/DM-26295"
if [ "$REPO_NAME" != "$OPERATOR" ]; then
  TMP_DIR=$(mktemp -d --suffix "$OPERATOR")/qserv-operator
  git clone --depth 1 -b "$GIT_REF" --single-branch https://github.com/lsst/"$OPERATOR".git "$TMP_DIR"
  "$TMP_DIR"/deploy.sh
else
  . "$DIR/env.sh"
  make install
  make deploy IMG="$OP_IMAGE"
  $DIR/tests/tools/wait-operator-ready.sh
fi
