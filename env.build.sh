VERSION="$(git describe --dirty --always)"

# Image version created by build procedure
OP_IMAGE="qserv/qserv-operator:$VERSION"

# Target namespace for ns-scoped operator
# use 'make yaml-ns-scoped' to create deployment files
NAMESPACE="qserv"
