#!/bin/ash

# Check wmgr is reachable on WORKER_POD, a Qserv worker pod
# WARN: this should be handled by wmgr readiness probe (TODO: check why probe check sometimes fails on CI)

set -e
set -x

usage() {
    cat << EOD

Usage: `basename $0` [WORKER_POD]

  Check wmgr is reachable on WORKER_POD, a Qserv worker pod

  WARN: Need to be run inside a pod

EOD
}

WMGR_PORT=5012

if [ $# -ge 2 ] ; then
    usage
    exit 2
elif [ $# -eq 1 ]; then
    WORKER_POD=$1
fi

while ! wget http://"$WORKER_POD":5012 2>&1 | \
    grep 'wget: server returned error: HTTP/1.0 401 UNAUTHORIZED' > /dev/null
do
    echo "Waiting for wmgr to be up on $WORKER_POD"
    sleep 2 
done
echo "wmgr ready on $WORKER_POD"