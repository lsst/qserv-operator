#!/usr/bin/env bash

# Push operator image to docker hub and produce related yaml file
# Publish a qserv-operator release

# @author  Fabrice Jammes, IN2P3

set -exuo pipefail

OP_VERSION=""
releasetag=""

DIR=$(cd "$(dirname "$0")"; pwd -P)

set -e

usage() {
  cat << EOD

Usage: `basename $0` [options] RELEASE_TAG

  Available options:
    -h          this message

Create a qserv-operator release tagged "RELEASE_TAG"
RELEASE_TAG must be of the form YYYY.M.D-rcX
- Push operator image to docker hub
- Produce operator.yaml and operator-ns-scoped.yaml
EOD
}

# get the options
while getopts ht: c ; do
    case $c in
	    h) usage ; exit 0 ;;
	    \?) usage ; exit 2 ;;
    esac
done
shift `expr $OPTIND - 1`

if [ $# -ne 1 ] ; then
    usage
    exit 2
fi

releasetag=$1
export OP_VERSION="$releasetag"

. "$DIR/env.build.sh"

make yaml yaml-ns-scoped
$DIR/build.sh
# Make file below compliant with goimport requirements
git checkout $DIR/api/v1alpha1/zz_generated.deepcopy.go

echo "Update Qserv images in manifests/base/image.yaml"
sed -ri  "s/^(\s*image: qserv\/.*:).*/\1$releasetag/" $DIR/manifests/base/image.yaml
git add .
git commit -m "Release $releasetag" || echo "Nothing to commit"
git tag -a "$releasetag" -m "Version $releasetag"
git push --tag
$DIR/push-image.sh
