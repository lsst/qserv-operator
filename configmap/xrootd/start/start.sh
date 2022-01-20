#!/bin/sh

# Start cmsd and xrootd inside pod
# Launch as qserv user

# @author  Fabrice Jammes, IN2P3

set -eux

usage() {
    cat << EOD

Usage: `basename $0` [options] [cmd]

  Available options:
    -S <service> Service to start, default to xrootd

  Prepare cmsd and xrootd (ulimit setup) startup and
  launch associated startup script using qserv user.
EOD
}

service=xrootd

# get the options
while getopts S: c ; do
    case $c in
        S) service="$OPTARG" ;;
        \?) usage ; exit 2 ;;
    esac
done
shift $(($OPTIND - 1))

if [ $# -ne 0 ] ; then
    usage
    exit 2
fi

XROOTD_RDR_DN="{{.XrootdRedirectorDn}}"

if hostname | egrep "^${XROOTD_RDR_DN}-[0-9]+$"
then
    COMPONENT_NAME='manager'
else
    COMPONENT_NAME='worker'
fi

CONFIG_DIR="/config-etc"
XROOTD_CONFIG="$CONFIG_DIR/xrootd.cf"
OPT_XRD_SSI=""

# COMPONENT_NAME is required by xrdssi plugin to
# choose which type of queries to launch against metadata
if [ "$COMPONENT_NAME" = 'worker' ]; then

    QSERV_WORKER_DB_USER="qsmaster"
    QSERV_WORKER_DB_PASSWORD=
    QSERV_WORKER_DB_DN="127.0.0.1"
    QSERV_WORKER_DB_PORT="3306"
    QSERV_WORKER_DB="qservw_worker"
    QSERV_WORKER_DB_URL="mysql://${QSERV_WORKER_DB_USER}:${QSERV_WORKER_DB_PASSWORD}@${QSERV_WORKER_DB_DN}:${QSERV_WORKER_DB_PORT}/${QSERV_WORKER_DB}"
    XRDSSI_CONFIG="$CONFIG_DIR/xrdssi.cf"

    # Wait for at least one xrootd redirector readiness
    until timeout 1 bash -c "cat < /dev/null > /dev/tcp/${XROOTD_RDR_DN}/2131"
    do
        echo "Wait for xrootd redirector to be up and running  (${XROOTD_RDR_DN})..."
        sleep 2
    done

    # xrootd/cmsd will use this configuration to learn the worker identity directly from
    # the worker's MySQL database. The plugin will automatically wait before
    # the database service will start up and be ready. So, it's no loner required to
    # track the availability and status of the database service by this script.
    export VNID_ARGS="${QSERV_WORKER_DB_URL} 0 0"

    OPT_XRD_SSI="-l @libXrdSsiLog.so -+xrdssi $XRDSSI_CONFIG"
fi

# Start service
#
echo "Start $service"
"$service" -c "$XROOTD_CONFIG" -n "$COMPONENT_NAME" -I v4 $OPT_XRD_SSI
