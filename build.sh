#!/usr/bin/env bash

# See
# https://github.com/operator-framework/operator-sdk/blob/master/doc/user-guide.md#build-and-run-the-operator

# @author  Fabrice Jammes, IN2P3

set -euxo pipefail

DIR=$(cd "$(dirname "$0")"; pwd -P)
. "$DIR/env.sh"



set -e

usage() {
  cat << EOD

Usage: `basename $0` [options] path host [host ...]

  Available options:
    -h          this message
    -k          development mode: load image in kind

Build qserv-operator image from source code.
EOD
}

kind=false

# get the options
while getopts hk c ; do
    case $c in
	    h) usage ; exit 0 ;;
	    k) kind=true ;;
	    \?) usage ; exit 2 ;;
    esac
done
shift `expr $OPTIND - 1`

if [ $# -ne 0 ] ; then
    usage
    exit 2
fi

make manifests
make docker-build IMG="$OP_IMAGE"

if [ $kind = true ]; then
  kind load docker-image "$OP_IMAGE"
fi
