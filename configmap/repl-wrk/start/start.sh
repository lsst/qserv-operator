#!/bin/bash

# Start Qserv replication worker service inside pod

# @author  Fabrice Jammes, IN2P3/SLAC

# WARN: password are displayed in debug logs
set -euxo pipefail

. /secret-mariadb/mariadb.secret.sh
. /secret-repl-db/repl-db.secret.sh

REPL_DB_PORT="3306"
REPL_DB_USER="qsreplica"
REPL_DB="qservReplica"
DATA_DIR="/qserv/data"
MYSQLD_DATA_DIR="$DATA_DIR/mysql"
MYSQLD_SOCKET="$MYSQLD_DATA_DIR/mysql.sock"
MYSQLD_USER_QSERV="qsmaster"
QSERV_WORKER_DB_DN="127.0.0.1"
QSERV_WORKER_DB_PORT="3306"
QSERV_WORKER_DB="qservw_worker"
QSERV_WORKER_DB_USER="root"
QSERV_WORKER_DB_PASSWORD=${MYSQL_ROOT_PASSWORD}


# Source pathes to eups packages
. /qserv/run/etc/sysconfig/qserv

WORKER_ID=$(hostname)

# Required by dataloader
mkdir -p "$DATA_DIR/ingest"

# Wait for remote repl-db started and contactable
while true; do
    if mysql --host="$REPL_DB_DN" --port="$REPL_DB_PORT" --user="$REPL_DB_USER" \
    --password="${MYSQL_REPLICA_PASSWORD}" --skip-column-names \
        "${REPL_DB}" -e "SELECT CONCAT('Mariadb is up: ', version())"
    then
        break
    else
        echo "Wait for repl-db startup"
    fi
    sleep 2
done

# Wait for all repl-wrk to be registered inside repl-db
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
        echo "Wait for all replication workers to register inside replication database"
    fi
    sleep 2
done

export LSST_LOG_CONFIG="/config-etc/log4cxx.replication.properties"

CONFIG="mysql://${REPL_DB_USER}:${MYSQL_REPLICA_PASSWORD}@${REPL_DB_DN}:${REPL_DB_PORT}/${REPL_DB}"
QSERV_WORKER_DB_URL="mysql://${QSERV_WORKER_DB_USER}:${QSERV_WORKER_DB_PASSWORD}@${QSERV_WORKER_DB_DN}:${QSERV_WORKER_DB_PORT}/${QSERV_WORKER_DB}"
qserv-replica-worker ${WORKER_ID} --config=${CONFIG} --qserv-worker-db="${QSERV_WORKER_DB_URL}" --debug

# For debug purpose
#while true;
#do
#    sleep 3600
#done
