#!/bin/sh

# Configure mariadb for qserv or replication service :
# - create data directory
# - create root password
# - create qserv/repl databases and user
# - deploy scisql plugin (qserv only)

# @author  Fabrice Jammes, IN2P3/SLAC

set -eu

# WARN: password are displayed in debug logs
# set -x

MYSQL_REPLICA_PASSWORD=''
MYSQL_MONITOR_PASSWORD=''

# Require root privileges
##
MARIADB_CONF="/config-etc/my.cnf"
if [ -e "$MARIADB_CONF" ]; then
    mkdir -p /etc/mysql
    ln -sf "$MARIADB_CONF" /etc/mysql/my.cnf
fi

if ! id 1000 > /dev/null 2>&1
then
    useradd qserv --uid 1000 --no-create-home
fi

if [ "$COMPONENT_NAME" = "repl" ]; then
    MYSQL_INSTALL_DB="mysql_install_db"
    . /secret-repl-db/repl-db.secret.sh
else
    # Source pathes to eups packages
    . /qserv/run/etc/sysconfig/qserv
    MYSQL_INSTALL_DB="${MYSQL_DIR}/scripts/mysql_install_db --basedir=$MYSQL_DIR"
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

SQL_DIR="/config-sql"

EXCLUDE_DIR1="lost+found"
DATA_FILES=$(find "$DATA_DIR" -mindepth 1 ! -name "$EXCLUDE_DIR1")

# Keep crashing if data initialization has failed in a previous instance
# of this script
STATE_FILE="$DATA_DIR/INIT_IN_PROGRESS.state"
if [ -f "$STATE_FILE" ]; then
    >&2 echo "ERROR: previous data initialization crashed"
    sleep 3600
    exit 1
fi

if [ ! "$DATA_FILES" ]
then
    touch "$STATE_FILE"
    echo "-- "
    echo "-- Installing mysql database files."
    ${MYSQL_INSTALL_DB} >/dev/null ||
        {
            echo "ERROR : mysql_install_db failed, exiting"
            exit 1
        }

    echo "-- "
    echo "-- Start mariadb server."
    mysqld &
    sleep 5

    echo "-- "
    echo "-- Change mariadb root password."
    mysqladmin -u root password "$MYSQL_ROOT_PASSWORD"

    echo "-- "
    echo "-- Initializing Qserv database"
    for file_name in "${SQL_DIR}/${COMPONENT_NAME}"/*; do
        echo "-- Loading ${file_name} in MySQL"
        sql_file_name="/tmp/out.sql"
        file_ext="${file_name#*\.}"
        if [ "${file_ext}" = "tpl.sql" ]; then
            awk \
                -v VAR1=${MYSQL_MONITOR_PASSWORD} \
                -v VAR2=${MYSQL_REPLICA_PASSWORD} \
                '{gsub(/<MYSQL_MONITOR_PASSWORD>/, VAR1);
                gsub(/<MYSQL_REPLICA_PASSWORD>/, VAR2);
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

    if [ "$COMPONENT_NAME" != "repl" ]; then
        echo "-- "
        echo "-- Deploy scisql plugin"
        # WARN: SciSQL shared library (libcisql*.so) deployed by command
        # below will be removed at each container startup.
        # That's why this shared library is currently
        # installed in mysql plugin directory at image creation.
        echo "$MYSQL_ROOT_PASSWORD" | scisql-deploy.py --mysql-dir="$MYSQL_DIR" \
            --mysql-socket="$MYSQLD_SOCKET"
    fi

    echo "-- Stop mariadb server."
    mysqladmin -u root --password="$MYSQL_ROOT_PASSWORD" shutdown
    rm "$STATE_FILE"
else
    echo "WARN: Skip mysqld initialization because of non empty $DATA_DIR:"
    ls -l "$DATA_DIR"
fi
