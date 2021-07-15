#!/usr/bin/env bash

# Push operator image to docker hub and produce related yaml file
# Publish a qserv-operator release

# @author  Fabrice Jammes, IN2P3

set -euo pipefail

OP_VERSION=""
releasetag=""

DIR=$(cd "$(dirname "$0")"; pwd -P)

set -e

usage() {
  cat << EOD

Usage: `basename $0` [options]

  Available options:
    -h          this message
    -t          release-tag: create a git release tag and use it to tag qserv-operator image

- Push operator image to docker hub
- Produce operator.yaml and operator-ns-scoped.yaml
EOD
}

# get the options
while getopts ht: c ; do
    case $c in
	    h) usage ; exit 0 ;;
      t) releasetag="$OPTARG" ;;
	    \?) usage ; exit 2 ;;
    esac
done
shift `expr $OPTIND - 1`

if [ -n "$releasetag" ] ; then
    export OP_VERSION="$releasetag"
fi

. "$DIR/env.build.sh"

if [ $# -ne 0 ] ; then
    usage
    exit 2
fi

make yaml yaml-ns-scoped
make docker-build IMG="$OP_IMAGE"
# WARN: Hack used to pass CI static code checks
git checkout $DIR/api/v1alpha1/zz_generated.deepcopy.go
docker push "$OP_IMAGE"

echo "-- WARNING Update Qserv images in manifests/base/image.yaml!!!"
echo "-- Then run command below to publish the release:"
echo "git add . &&  git commit -m "Release $releasetag" && git tag -a "$releasetag" -m "Version $releasetag" && git push --tag"
echo "-- Rebuild and push Qserv image with release tag:"
echo "./build.sh"
echo "./push-image.sh -d"
