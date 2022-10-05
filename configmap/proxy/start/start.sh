#!/bin/sh
#
# mysql-proxy This script starts mysql-proxy
#
# description: mysql-proxy is a proxy daemon for mysql

# Description: mysql-proxy is the user (i.e. mysql-client) frontend for \
#              Qserv czar service. \
#              It receive SQL queries, process it using lua plugin, \
#              and send it to Qserv czar. \
#              Once Qserv czar have returned the results, mysql-proxy \
#              sends it to mysql-client. \

set -ex

PASSWORD=CHANGEME

# Hack to upgrade database schema
# FIXME improve it
entrypoint --log-level DEBUG smig-update --czar-connection mysql://root:$PASSWORD@localhost:3306

entrypoint --log-level DEBUG proxy \
  --db-uri "mysql://qsmaster@127.0.0.1:3306?socket={{.MariadbSocket}}" \
  --db-admin-uri "mysql://root:$PASSWORD@127.0.0.1:3306?socket={{.MariadbSocket}}" \
  --xrootd-manager qserv-xrootd-redirector \
  --log-cfg-file /cm-etc/log.cnf
