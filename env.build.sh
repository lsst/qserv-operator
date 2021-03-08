VERSION="$(git describe --dirty --always)"

# Image version created by build procedure
OP_IMAGE="qserv/qserv-operator:$VERSION"
