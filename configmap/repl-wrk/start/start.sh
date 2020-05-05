#!/bin/sh

# Start Qserv replication worker service inside pod

# @author  Fabrice Jammes, IN2P3/SLAC

set -e
# WARN: password are displayed in debug logs
set -x

REPL_DB_PORT="3306"
REPL_DB_USER="qsreplica"
REPL_DB="qservReplica"
DATA_DIR="/qserv/data"
MYSQLD_DATA_DIR="$DATA_DIR/mysql"
MYSQLD_SOCKET="$MYSQLD_DATA_DIR/mysql.sock"
MYSQLD_USER_QSERV="qsmaster"

. /secret-mariadb/mariadb.secret.sh
. /secret-repl-db/repl-db.secret.sh

# Add mysql client to path
export PATH="/stack/stack/current/Linux64/mariadb/10.2.14.lsst3-1-g07c67f4/bin/:$PATH"

WORKER_ID=$(hostname)


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
qserv-replica-worker ${WORKER_ID} --config=${CONFIG} --qserv-db-password="${MYSQL_ROOT_PASSWORD}" --debug

# For debug purpose
#while true;
#do
#    sleep 3600
#done
