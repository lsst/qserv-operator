RELEASE=0.0.1
GIT_HASH="$(git rev-parse --short HEAD)"
VERSION="$RELEASE-${GIT_HASH}"

# Image version created by build procedure
OP_IMAGE="qserv/qserv-operator:$VERSION"
