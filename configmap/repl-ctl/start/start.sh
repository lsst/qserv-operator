#!/bin/sh

# Start Qserv replication controller service inside pod

# @author  Fabrice Jammes, IN2P3/SLAC

set -e
# WARN: password are displayed in debug logs
set -x

# Load parameters of the setup into the corresponding environment
# variables

# Start master controller

REPL_DB_PORT="3306"
REPL_DB_USER="qsreplica"
REPL_DB="qservReplica"

# Add mysql client to path
export PATH="/stack/stack/current/Linux64/mariadb/10.2.14.lsst3-1-g07c67f4/bin/:$PATH"

. /secret-mariadb/mariadb.secret.sh
. /secret-repl-db/repl-db.secret.sh

echo "Start replication controller pod: ${HOSTNAME}"

# Wait for repl-db started
# and contactable
while true; do
    if mysql --host="$REPL_DB_DN" --port="$REPL_DB_PORT" --user="$REPL_DB_USER" \
      --password="$MYSQL_REPLICA_PASSWORD" --skip-column-names \
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
    REGISTERED_WORKERS=$(mysql --host="$REPL_DB_DN" --port="$REPL_DB_PORT" \
    --user="$REPL_DB_USER" --password="$MYSQL_REPLICA_PASSWORD" \
    --skip-column-names --batch "${REPL_DB}" -e "SELECT count(*) from config_worker")
    if [ "$REGISTERED_WORKERS" -eq "$WORKER_COUNT" ]
    then
        echo "Replication workers all registered inside replication database: \
        (${REGISTERED_WORKERS}/${WORKER_COUNT})"
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

export LSST_LOG_CONFIG="/config-etc/log4cxx.replication.properties"

CONFIG="mysql://${REPL_DB_USER}:${MYSQL_REPLICA_PASSWORD}@${REPL_DB_DN}:${REPL_DB_PORT}/${REPL_DB}"
PARAMETERS="--worker-evict-timeout=3600 --health-probe-interval=120 --replication-interval=1200"
MALLOC_CONF=${OPT_MALLOC_CONF} LD_PRELOAD=${OPT_LD_PRELOAD} \
qserv-replica-master-http ${PARAMETERS} --config="${CONFIG}" --qserv-db-password="${MYSQL_ROOT_PASSWORD}"

# For debug purpose
#while true;
#do
#    sleep 3600
#done
