#!/bin/bash

set -euxo pipefail

DIR=$(cd "$(dirname "$0")"; pwd -P)

usage() {
    cat << EOD

Usage: `basename $0` [options]

  Available options:
    -n <namespace>  Target namespace
    -h              Tthis message

  Install qserv-operator inside a given namespace

EOD
}

NAMESPACE_OPT=""
NAMESPACE=""

# get the options
while getopts n:h c ; do
    case $c in
        n) NAMESPACE="$OPTARG" ;;
        h) usage ; exit 0 ;;
    esac
done
shift $(($OPTIND - 1))

if [ $# -ge 1 ] ; then
    usage
    exit 2
fi

if [ "$NAMESPACE" ]; then
    NAMESPACE_OPT="--namespace=$NAMESPACE"
fi

kapply="kubectl apply $NAMESPACE_OPT -f "

$kapply "$DIR"/deploy/crds/qserv_v1alpha1_qserv_crd.yaml
$kapply "$DIR"/deploy/service_account.yaml
$kapply "$DIR"/deploy/role.yaml
$kapply "$DIR"/deploy/role_binding.yaml
$kapply "$DIR"/deploy/operator.yaml

echo "----------------------------------"
echo "Run command below to deploy Qserv:"
echo "----------------------------------"
echo "kubectl apply -k $DIR/base $NAMESPACE_OPT"
