#!/bin/bash
# Tag a qserv version
# use on master branch only

set -euxo pipefail

VERSION="v$(date +'%Y.%m').1"
git tag -a "$VERSION" -m "Version $VERSION"
