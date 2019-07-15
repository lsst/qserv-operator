#!/bin/ sh

# Start Qserv replication worker service inside pod

# @author  Fabrice Jammes, IN2P3/SLAC

set -e
set -x

REPL_DB_HOST="repl-db-0.qserv"
REPL_DB_PORT="3306"
REPL_DB_USER="qsreplica"
REPL_DB="qservReplica"
DATA_DIR="/qserv/data"
MYSQLD_DATA_DIR="$DATA_DIR/mysql"
MYSQLD_SOCKET="$MYSQLD_DATA_DIR/mysql.sock"
MYSQLD_USER_QSERV="qsmaster"


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

 # Retrieve worker id
WORKER_ID=$(mysql --socket "$MYSQLD_SOCKET" --batch \
    --skip-column-names --user="$MYSQLD_USER_QSERV" -e "SELECT id FROM qservw_worker.Id;")
if [ -z "$WORKER_ID" ]; then
    >&2 echo "ERROR: unable to retrieve worker id for $HOSTNAME"
    exit 1 
fi

HOST_DN=$(hostname --fqdn)

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

# Register repl-wrk on repl-db
SQL="INSERT INTO \`config_worker\` VALUES ('${WORKER_ID}', 1, 0, '${HOST_DN}', \
    NULL, '${HOST_DN}',  NULL, NULL) ON DUPLICATE KEY UPDATE name='${WORKER_ID}', \
    svc_host='${HOST_DN}', fs_host='${HOST_DN}';"
mysql --host="$REPL_DB_HOST" --port="$REPL_DB_PORT" --user="$REPL_DB_USER" -vv \
    "${REPL_DB}" -e "$SQL"

LSST_LOG_CONFIG="/config-etc/log4cxx.replication.properties"

CONFIG="mysql://${REPL_DB_USER}@${REPL_DB_HOST}:${REPL_DB_PORT}/${REPL_DB}"
qserv-replica-worker ${WORKER_ID} --config=${CONFIG}

# For debug purpose
#while true;
#do
#    sleep 5
#done
