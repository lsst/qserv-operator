#!/bin/bash

# Install and uninstall Qserv operator

set -euxo pipefail

usage() {
    cat << EOD

Usage: `basename $0` [options] [cmd]

  Available options:
    -h, --help           This message
    -d, --dev            Install from local git repository
    --n NAMESPACE        Specify namespace (default: kube-system)
    --install-kubedb     Install KubeDB operator
    --uninstall          Uninstall Qserv-operator,
                         and related CustomResourceDefinition/CustomResource

  Install Qserv-operator, and eventually KubeDB-operator.

EOD
}

DIR=$(cd "$(dirname "$0")"; pwd -P)

echo "Check kubeconfig context"
kubectl config current-context || {
  echo "Set a context (kubectl use-context <context>) out of the following:"
  echo
  kubectl config get-contexts
  exit 1
}
echo ""

KUBEDB=false
KUBEDB_URL="https://raw.githubusercontent.com/kubedb/installer/89fab34cf2f5d9e0bcc3c2d5b0f0599f94ff0dca/deploy/kubedb.sh"

DEV_INSTALL=false
UNINSTALL=false
PURGE=true

NAMESPACE=$(kubectl config view --minify --output 'jsonpath={..namespace}')
NAMESPACE=${NAMESPACE:-default}
export QSERV_DOCKER_REGISTRY=${QSERV_DOCKER_REGISTRY:-qserv}
export QSERV_OPERATOR_TAG=${QSERV_OPERATOR_TAG:-v0.13.0-rc.0}


while test $# -gt 0; do
  case "$1" in
    -h | --help)
      usage
      exit 0
      ;;
    -d | --dev)
      DEV_INSTALL=true
      shift
      ;;
    -n)
      shift
      NAMESPACE="$1"
      shift
      ;;
    --install-kubedb)
      KUBEDB=true
      shift
      ;;
      --uninstall)
      export UNINSTALL=true
      shift
      ;;
    --purge)
      export PURGE=true
      shift
      ;;
    *)
      echo "Error: unknown flag:" $1
      usage
      exit 1
      ;;
  esac
done

if [ "$UNINSTALL" = true ]; then

  kubectl delete deployment,role,rolebinding,serviceaccount qserv-operator
  kubectl delete crds qservs.qserv.lsst.org

  (
  # Put KubeDB backup files (yaml CRDs) inside a temporary directory
  TMP_DIR=$(mktemp -d --suffix="-qserv-operator")
  cd "$TMP_DIR"
  curl -fsSL "$KUBEDB_URL" | bash -s -- -n "$NAMESPACE" --uninstall --purge
  )
  echo
  echo "Successfully uninstall Qserv-operator"
  exit 0
fi

if [ "$DEV_INSTALL" = true ]; then
  MANIFESTS_DIR=$(dirname "$DIR")
else
  MANIFESTS_DIR="https://raw.githubusercontent.com/lsst/qserv-operator/master"
fi

kapply="kubectl apply -n $NAMESPACE -f "

echo "Install Qserv-operator"

$kapply "$MANIFESTS_DIR"/deploy/crds/qserv_v1alpha1_qserv_crd.yaml
$kapply "$MANIFESTS_DIR"/deploy/service_account.yaml
$kapply "$MANIFESTS_DIR"/deploy/role.yaml
$kapply "$MANIFESTS_DIR"/deploy/role_binding.yaml
$kapply "$MANIFESTS_DIR"/deploy/operator.yaml

if [ "$KUBEDB" = true ]; then
  (
  # Put KubeDB install files inside a temporary directory
  TMP_DIR=$(mktemp -d --suffix="-qserv-operator")
  cd "$TMP_DIR"
  echo "Install KubeDB"
  # See https://kubedb.com/docs/v0.13.0-rc.0/setup/install/, but installer is broken
  curl -fsSL "$KUBEDB_URL" | KUBEDB_CATALOG=redis bash -s -- -n "$NAMESPACE"
  kubectl wait --for=condition=Ready pods -l app=kubedb -n "$NAMESPACE"
  )
fi

while ! kubectl wait --for=condition=Ready pods -l name=qserv-operator -n "$NAMESPACE"
do
  echo "Waiting for operator to be ready..."
  kubectl describe pod -l name=qserv-operator -n "$NAMESPACE"
done

echo
echo "Successfully installed Qserv operator in '$NAMESPACE' namespace."
