#!/usr/bin/env bash

# Push operator image to docker hub and produce related yaml file
# Publish a qserv-operator release

# @author  Fabrice Jammes, IN2P3

set -exuo pipefail

OP_VERSION=""
releasetag=""

DIR=$(cd "$(dirname "$0")"; pwd -P)

usage() {
  cat << EOD

Usage: `basename $0` [options] RELEASE_TAG

  Available options:
    -h          this message

Create a qserv-operator release tagged "RELEASE_TAG"
- Release tag template YYYY.M.<i>-rc<j>, i and j are integers
- Create a git release tag and use it to tag qserv-operator image
- Push operator image to docker hub
- Produce operator.yaml and operator-ns-scoped.yaml
- Produce operatorHub bundle in bundle/ directory
EOD
}

# get the options
while getopts h c ; do
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

releasetag="$1"
export OP_VERSION="$releasetag"

. "$DIR/env.build.sh"

make yaml yaml-ns-scoped
$DIR/build.sh
# Make file below compliant with goimport requirements
git checkout $DIR/api/v1alpha1/zz_generated.deepcopy.go

echo "Update Qserv images in manifests/base/image.yaml"
sed -ri  "s/^(\s*image: qserv\/.*:).*/\1$releasetag/" $DIR/manifests/base/image.yaml
echo "Update release number in documentation"
find $DIR/doc -type f -print0 | xargs -0 sed -ri  "s/RELEASE=\".*\"/RELEASE=\"$releasetag\"/"
sed -ri  "s/RELEASE=\".*\"/RELEASE=\"$releasetag\"/" $DIR/README.md

# Prepare operatorHub files
# Edit 'replaces', 'image' and 'containerImage' fields in config/manifests/bases/qserv-operator.clusterserviceversion.yaml
csv_file="$DIR/config/manifests/bases/qserv-operator.clusterserviceversion.yaml"
previous_version=$(grep -oP 'qserv\/qserv-operator:([0-9]+\.[0-9]+\.[0-9](-rc[0-9]+)?)' "$csv_file" | cut -d: -f2)
sed -ri  "s/replaces: qserv-operator\.v([0-9]+\.[0-9]+\.[0-9](-rc[0-9]+)?)/replaces: qserv-operator\.v$previous_version/"  "$csv_file"
sed -ri  "s/qserv\/([a-z\-]+):([0-9]+\.[0-9]+\.[0-9](-rc[0-9]+)?)/qserv\/\1:$releasetag/" "$csv_file"
# TODO Replace §RELEASE env value in CSV documentation

git add .
git commit -m "Release $releasetag" || echo "Nothing to commit"
git tag -a "$releasetag" -m "Version $releasetag"
git push --tag
$DIR/push-image.sh

echo "Publish release to operator hub"
make bundle
