#!/bin/sh

# Start cmsd and xrootd inside pod
# Launch as qserv user

# @author  Fabrice Jammes, IN2P3

set -eux

# Increase limit for locked-in-memory size
MLOCK_AMOUNT=$(grep MemTotal /proc/meminfo | awk '{printf("%.0f\n", $2 - 1000000)}')
ulimit -l "$MLOCK_AMOUNT"

# FIXME password is required for database initialization and this should move to a dedicated container
su qserv -c 'MYSQL_ROOT_PASSWORD="CHANGEME" && \
entrypoint --log-level DEBUG worker-xrootd \
          --db-uri "{{.SocketQservUser}}" \
          --db-admin-uri "{{.SocketRootUser}}" \
          --vnid-config "@/usr/local/lib64/libreplica.so {{.WorkerDatabaseLocalURL}} 0 0" \
          --cmsd-manager-name "{{.XrootdRedirectorDN}}" \
          --cmsd-manager-count "{{.XrootdRedirectorReplicas}}" \
          --mysql-monitor-password CHANGEME_MONITOR'
