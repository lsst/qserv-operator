#!/bin/sh

# Manage GAIA data

set -e
set -x

usage() {
    cat << EOD
Usage: $(basename "$0") [options]
Available options:
  -h            This message

Init k8s master

EOD
}

# Get the options
while getopts h c ; do
    case $c in
        h) usage ; exit 0 ;;
        \?) usage ; exit 2 ;;
    esac
done
shift "$((OPTIND-1))"

if [ $# -ne 0 ] ; then
    usage
    exit 2
fi

DIR=$(cd "$(dirname "$0")"; pwd -P)

NODES=$(kubectl get nodes -o=jsonpath='{range .items[*]}{.metadata.name}{"\n"}{end}' | grep qserv | grep -v dax | tr '\n' ' ')

echo "Copy scripts to all nodes"
echo "-------------------------"
parallel -vvv --tag -- "scp -r $DIR/resource $USER@{}:/tmp" ::: $NODES

echo "Move data"
echo "---------"
parallel -vvv --tag -- "ssh $USER@{} -- sudo 'sh /tmp/resource/mv-data.sh'" ::: $NODES