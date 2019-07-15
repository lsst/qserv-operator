#!/bin/sh

# Start Qserv replication controller service inside pod

# @author  Fabrice Jammes, IN2P3/SLAC

set -e
set -x

# Load parameters of the setup into the corresponding environment
# variables

# Start master controller

REPL_DB_HOST="repl-db-0.qserv"
REPL_DB_PORT="3306"
REPL_DB_USER="qsreplica"
REPL_DB="qservReplica"

echo "Start replication controller pod: ${HOSTNAME}"

# Wait for repl-db started
# and contactable
while true; do
    if mysql --host="$REPL_DB_HOST" --port="$REPL_DB_PORT" --user="$REPL_DB_USER" --skip-column-names \
        "${REPL_DB}" -e "SELECT CONCAT('Mariadb is up: ', version())"
    then
        break
    else
        echo "Wait for repl-db startup"
    fi
    sleep 2
done


# Wait for repl-wrk to register inside repl-db
while true; do
    REGISTERED_WORKERS=$(mysql --host="$REPL_DB_HOST" --port="$REPL_DB_PORT" --user="$REPL_DB_USER" \
    --skip-column-names --batch "${REPL_DB}" -e "SELECT count(*) from config_worker")
    if [ "$REGISTERED_WORKERS" -eq "$WORKER_COUNT" ]
    then
        break
    else
        echo "Wait for all replication workers to register inside replication database: \
        (${REGISTERED_WORKERS}/${WORKER_COUNT})"
    fi
    sleep 2
done

OPT_MALLOC_CONF=
OPT_LD_PRELOAD=
if [ ! -z "${USE_JEMALLOC}" ]; then
    OPT_MALLOC_CONF=prof_leak:true,lg_prof_interval:31,lg_prof_sample:22,prof_final:true
    OPT_LD_PRELOAD=/qserv/lib/libjemalloc.so
fi

# Work directory for the applications. It can be used by applications
# to store core files, as well as various debug information.
# TODO: Enable core dump management inside Kubernetes
WORK_DIR="/tmp"
cd "${WORK_DIR}"

LSST_LOG_CONFIG="/config-etc/log4cxx.replication.properties"

CONFIG="mysql://${REPL_DB_USER}@${REPL_DB_HOST}:${REPL_DB_PORT}/${REPL_DB}"
PARAMETERS="--worker-evict-timeout=3600 --health-probe-interval=120 --replication-interval=1200"
MALLOC_CONF=${OPT_MALLOC_CONF} LD_PRELOAD=${OPT_LD_PRELOAD} \
qserv-replica-master-http ${PARAMETERS} --config="${CONFIG}"

# For debug purpose
#while true;
#do
#    sleep 5
#done
