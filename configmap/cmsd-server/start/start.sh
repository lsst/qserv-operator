#!/bin/sh

# Start cmsd and xrootd inside pod
# Launch as qserv user

# @author  Fabrice Jammes, IN2P3

set -eux

entrypoint --log-level DEBUG worker-cmsd \
          --db-uri {{.SocketQservUser}} \
          --vnid-config "@/usr/local/lib64/libreplica.so {{.WorkerDatabaseLocalURL}} 0 0" \
          --cmsd-manager-name {{.XrootdRedirectorDN}} \
          --cmsd-manager-count {{.XrootdRedirectorReplicas}}
