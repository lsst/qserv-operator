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

<<<<<<< HEAD
Usage: `basename $0` [options] RELEASE_TAG
=======
Usage: `basename $0` [options] release-tag
>>>>>>> 2b1aa0a... Release 2021.6.3-rc2

  Available options:
    -h          this message

<<<<<<< HEAD
Create a qserv-operator release tagged "RELEASE_TAG"
=======
>>>>>>> 2b1aa0a... Release 2021.6.3-rc2
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
make bundle
# Make file below compliant with goimport requirements
git checkout $DIR/api/v1alpha1/zz_generated.deepcopy.go

echo "Update Qserv images in manifests/base/image.yaml"
sed -ri  "s/^(\s*image: qserv\/.*:).*/\1$releasetag/" $DIR/manifests/base/image.yaml
echo "Update release number in documentation"
find $DIR/doc -type f -print0 | xargs -0 sed -ri  "s/RELEASE=\".*\"/RELEASE=\"$releasetag\"/"
sed -ri  "s/RELEASE=\".*\"/RELEASE=\"$releasetag\"/" $DIR/README.md
git add .
git commit -m "Release $releasetag" || echo "Nothing to commit"
git tag -a "$releasetag" -m "Version $releasetag"
git push --tag
$DIR/push-image.sh

