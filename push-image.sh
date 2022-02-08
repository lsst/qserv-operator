#!/usr/bin/env bash

# Push image to Docker Hub or load it inside kind

# @author  Fabrice Jammes, IN2P3

set -euxo pipefail

DIR=$(cd "$(dirname "$0")"; pwd -P)
. "$DIR/env.build.sh"

IMAGE="$OP_IMAGE"

set -e

usage() {
  cat << EOD

Usage: `basename $0` [options] path host [host ...]

  Available options:
    -h          this message
    -k          development mode: load image in kind
    -d          push image to docker hub (default)

Push image to Docker Hub and/or load it inside kind
EOD
}

kind=false
dockerhub=false

# get the options
while getopts dhk c ; do
    case $c in
	    h) usage ; exit 0 ;;
	    k) kind=true ;;
	    d) dockerhub=true ;;
	    \?) usage ; exit 2 ;;
    esac
done
shift `expr $OPTIND - 1`

# Default to dockerhub
if    [ $kind = false ] && [ $dockerhub = false ]
then
    dockerhub=true
fi;


if [ $# -ne 0 ] ; then
    usage
    exit 2
fi

if [ $kind = true ]; then
  kind load docker-image "$IMAGE"
fi
if [ $dockerhub = true ]; then
  docker push "$IMAGE"
fi
