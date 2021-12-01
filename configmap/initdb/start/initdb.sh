#!/bin/bash

# Configure mariadb for qserv or replication service :
# - create data directory
# - create root password
# - create qserv/repl databases and user
# - deploy scisql UDF in MySQL database (only for qserv czar(s)/worker(s))

# @author  Fabrice Jammes, IN2P3/SLAC

set -euxo pipefail

# WARN: password are displayed in debug logs
# set -x

MYSQL_INGEST_PASSWORD=''
MYSQL_REPLICA_PASSWORD=''
MYSQL_MONITOR_PASSWORD=''
SQL_DIR="/config-sql"

# Used for Qserv czar and worker databases

if [ ! "$COMPONENT_NAME" = "czar" ] && [ ! "$COMPONENT_NAME" = "worker" ]; then
    . /secret-"$COMPONENT_NAME"/"$COMPONENT_NAME".secret.sh
    INSTALL_SCISQL=false
else
    # Initialize scisql on both czar and worker
    INSTALL_SCISQL=true
fi

DATA_DIR="/qserv/data"
MYSQLD_DATA_DIR="$DATA_DIR/mysql"
MYSQLD_SOCKET="$MYSQLD_DATA_DIR/mysql.sock"

# Load mariadb secrets
. /secret-mariadb/mariadb.secret.sh
if [ -z "$MYSQL_ROOT_PASSWORD" ]; then
    echo "ERROR : mariadb root password is missing, exiting"
    exit 2
fi


# Keep crashing if data initialization has failed in a previous instance
# of this script
STATE_FILE="$DATA_DIR/INIT_IN_PROGRESS.state"
if [ -f "$STATE_FILE" ]; then
    >&2 echo "ERROR: previous data initialization crashed"
    sleep 3600
    exit 1
fi

HOST="$(hostname)"

if [ ! -e "$MYSQLD_DATA_DIR" ]
then
    touch "$STATE_FILE"
    echo "-- "
    echo "-- Installing mysql database files."
    id
    ls -rtl /qserv/data
    mysql_install_db --auth-root-authentication-method=normal >/dev/null ||
        {
            echo "ERROR : mysql_install_db failed, exiting"
            exit 1
        }

    echo "-- "
    echo "-- Start mariadb server."
    # Skip networking so to prevent replication controller and workers startup
    mysqld --skip-networking &
    sleep 5

    echo "-- "
    echo "-- Change mariadb root password."
    mysqladmin -u root password "$MYSQL_ROOT_PASSWORD"

    echo "-- "
    echo "-- Initialize Qserv database"
    for file_name in "${SQL_DIR}/${COMPONENT_NAME}"/*; do
        echo "-- Load ${file_name} in MySQL"
        sql_file_name="/tmp/out.sql"
        file_ext="${file_name#*\.}"
        if [ "${file_ext}" = "tpl.sql" ]; then
            awk \
                -v INGEST_PASS=${MYSQL_INGEST_PASSWORD} \
                -v MON_PASS=${MYSQL_MONITOR_PASSWORD} \
                -v REPL_PASS=${MYSQL_REPLICA_PASSWORD} \
                -v HOST=${HOST} \
                '{gsub(/<MYSQL_INGEST_PASSWORD>/, INGEST_PASS);
                gsub(/<MYSQL_MONITOR_PASSWORD>/, MON_PASS);
                gsub(/<MYSQL_REPLICA_PASSWORD>/, REPL_PASS);
                gsub(/<HOST>/, HOST);
                print}' "$file_name" > "$sql_file_name"
        else
            sql_file_name="$file_name"
        fi
        if mysql -vvv --user="root" --password="${MYSQL_ROOT_PASSWORD}" \
            < "${sql_file_name}"
        then
            echo "-- -> success"
        else
            >&2 echo "-- -> error"
            exit 1
        fi
    done

    if [ "$INSTALL_SCISQL" = true ]; then
        if mysql -vvv --user="root" --password="${MYSQL_ROOT_PASSWORD}" \
            < "/docker-entrypoint-initdb.d/scisql.sql"
        then
            echo "-- -> success"
        else
            >&2 echo "-- -> error"
            exit 1
        fi
    fi

    echo "-- Stop mariadb server."
    mysqladmin -u root --password="$MYSQL_ROOT_PASSWORD" shutdown
    rm "$STATE_FILE"
else
    echo "WARN: Skip mysqld initialization because of non empty $DATA_DIR:"
    ls -l "$DATA_DIR"
fi
