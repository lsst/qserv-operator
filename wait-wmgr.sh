#!/bin/ash

# Check wmgr is reachable on WORKER_POD, a Qserv worker pod
# WARN: this should be handled by wmgr readiness probe (TODO: check why probe check sometimes fails on CI)

set -uxo pipefail

usage() {
    cat << EOD

Usage: `basename $0` [WORKER_POD]

  Check wmgr is reachable on WORKER_POD, a Qserv worker pod

  WARN: Need to be run inside a pod

EOD
}

WORKER_POD=""
WMGR_PORT=5012

if [ $# -ne 1 ] ; then
    usage
    exit 2
else
    WORKER_POD=$1
fi

echo "Wait for wmgr on $WORKER_POD"
READY=false
while [ $READY = false ]
do
    curl --output /dev/null --silent --head --fail http://$WORKER_POD:$WMGR_PORT
    if [ $? == 22 ]
    then
        READY=true
    else
        printf '.'
        sleep 5
    fi
done

echo "wmgr ready on $WORKER_POD"
