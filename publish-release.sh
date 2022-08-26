#!/usr/bin/env bash

# Push operator image to docker hub and produce related yaml file
# Publish a qserv-operator release

# @author  Fabrice Jammes, IN2P3

set -exuo pipefail

release=true
OP_VERSION=""
releasetag=""

DIR=$(cd "$(dirname "$0")"; pwd -P)

usage() {
  cat << EOD

Usage: `basename $0` [options] [RELEASE_TAG]

  Available options:
    -h          this message
    -m          minimal release, with not tag, used to publish from a ticket branch

Create a qserv-operator release tagged "RELEASE_TAG"
- Release tag template YYYY.M.<i>-rc<j>, i and j are integers
- Create a git release tag and use it to tag qserv-operator image
- Push operator image to docker hub
- Produce operator.yaml and operator-ns-scoped.yaml
- Produce operatorHub bundle in bundle/ directory
EOD
}

# get the options
while getopts hm c ; do
    case $c in
      h) usage ; exit 0 ;;
      m) release=false ;;
      \?) usage ; exit 2 ;;
    esac
done
shift `expr $OPTIND - 1`

if [ $# -gt 1 ] ; then
    usage
    exit 2
fi


if [ $# -eq 1 ] ; then
  releasetag="$1"
  export OP_VERSION="$1"
  message="Publish release"
else
  message="Publish version"
fi

if [[ "$releasetag" =~ '/' || "$releasetag" =~ '\' ]]
then
  >&2 echo "ERROR: Found '\' or '/' in release tag $releasetag"
  exit 1
fi

. "$DIR/env.build.sh"

$DIR/build.sh
$DIR/push-image.sh
make yaml yaml-ns-scoped

if [ "$release" = true ]; then

  echo "Update Qserv images in manifests/base/image.yaml"
  sed -ri  "s/^(\s*image: qserv\/.*:).*/\1$releasetag/" $DIR/manifests/base/image.yaml
  echo "Update release number in documentation"
  find $DIR/doc -type f -print0 | xargs -0 sed -ri  "s/RELEASE=\".*\"/RELEASE=\"$releasetag\"/"
  sed -ri  "s/RELEASE=\".*\"/RELEASE=\"$releasetag\"/" $DIR/README.md

  # Prepare operatorHub files
  # Edit 'replaces', 'image' and 'containerImage' fields in config/manifests/bases/qserv-operator.clusterserviceversion.yaml
  csv_file="$DIR/config/manifests/bases/qserv-operator.clusterserviceversion.yaml"
  sample_file="$DIR/config/samples/qserv_v1beta1_qserv.yaml"
  previous_version=$(grep -oP 'qserv\/qserv-operator:([0-9]+\.[0-9]+\.[0-9](-rc[0-9]+)?)' "$csv_file" | cut -d: -f2)
  sed -ri  "s/(202[0-9]+\.[0-9]+\.[0-9](-rc[0-9]+)?)/$releasetag/" "$sample_file" "$csv_file"
  sed -ri  "s/replaces: qserv-operator\.v([0-9]+\.[0-9]+\.[0-9](-rc[0-9]+)?)/replaces: qserv-operator\.v$previous_version/"  "$csv_file"

fi

git add .
git commit -m "$message $VERSION" || echo "Nothing to commit"

if [ "$release" = true ]; then
  git tag -a "$releasetag" -m "Version $releasetag"
  git push --follow-tags
fi

