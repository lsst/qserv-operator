GIT_HASH="$(git describe --dirty --always)"
TAG=${OP_VERSION:-${GIT_HASH}}

# Image version create by build procedure
OP_IMAGE="qserv/qserv-operator:$TAG"
