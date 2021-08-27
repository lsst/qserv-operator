GIT_HASH="$(git describe --dirty --always)"
VERSION=${OP_VERSION:-${GIT_HASH}}

# Image version created by build procedure
OP_IMAGE=${OP_IMAGE:-"qserv/qserv-operator:$VERSION"}

# Target namespace for ns-scoped operator
# use 'make yaml-ns-scoped' to create deployment files
NAMESPACE="qserv"
