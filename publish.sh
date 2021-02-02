#!/usr/bin/env bash

# Push operator image to docker hub and produce related yaml file 

# @author  Fabrice Jammes, IN2P3

set -euxo pipefail

DIR=$(cd "$(dirname "$0")"; pwd -P)
. "$DIR/env.build.sh"

set -e

usage() {
  cat << EOD

Usage: `basename $0` [options] path host [host ...]

  Available options:
    -h          this message

Push operator image to docker hub and produce related yaml file 
EOD
}

kind=false

# get the options
while getopts hk c ; do
    case $c in
	    h) usage ; exit 0 ;;
	    \?) usage ; exit 2 ;;
    esac
done
shift `expr $OPTIND - 1`

if [ $# -ne 0 ] ; then
    usage
    exit 2
fi

make yaml
make docker-build IMG="$OP_IMAGE"
# WARN: Hack used to pass CI static code checks
git checkout $DIR/api/v1alpha1/zz_generated.deepcopy.go
docker push "$OP_IMAGE"
