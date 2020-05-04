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

# Wait for local mysql to be started
while true; do
    if mysql --socket "$MYSQLD_SOCKET" --user="$MYSQLD_USER_QSERV"  --skip-column-names \
        -e "SELECT CONCAT('Mariadb is up: ', version())"
    then
        break
    else
        echo "Wait for MySQL startup"
    fi
    sleep 2
done

 # Retrieve worker id on local mysql
WORKER_ID=$(mysql --socket "$MYSQLD_SOCKET" --batch \
    --skip-column-names --user="$MYSQLD_USER_QSERV" -e "SELECT id FROM qservw_worker.Id;")
if [ -z "$WORKER_ID" ]; then
    >&2 echo "ERROR: unable to retrieve worker id for $HOSTNAME"
    exit 1
fi

HOST_DN=$(hostname --fqdn)

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

# Register repl-wrk on repl-db
SQL="INSERT INTO \`config_worker\` VALUES ('${WORKER_ID}', 1, 0, '${HOST_DN}', \
    NULL, '${HOST_DN}',  NULL, NULL, 'localhost', NULL, NULL, '${HOST_DN}', NULL, NULL, '${HOST_DN}', NULL, NULL) ON DUPLICATE KEY UPDATE name='${WORKER_ID}', \
    svc_host='${HOST_DN}', fs_host='${HOST_DN}', loader_host='${HOST_DN}', exporter_host='${HOST_DN}';"
mysql --host="$REPL_DB_DN" --port="$REPL_DB_PORT" --user="$REPL_DB_USER" \
--password="${MYSQL_REPLICA_PASSWORD}" -vv "${REPL_DB}" -e "$SQL"

export LSST_LOG_CONFIG="/config-etc/log4cxx.replication.properties"

CONFIG="mysql://${REPL_DB_USER}:${MYSQL_REPLICA_PASSWORD}@${REPL_DB_DN}:${REPL_DB_PORT}/${REPL_DB}"
qserv-replica-worker ${WORKER_ID} --config=${CONFIG} --qserv-db-password="${MYSQL_ROOT_PASSWORD}" --debug

# For debug purpose
#while true;
#do
#    sleep 3600
#done
