#!/bin/bash

# Start Qserv replication controller
# @author  Fabrice Jammes, IN2P3/SLAC

set -exo pipefail

. /secret-mariadb/mariadb.secret.sh
. /secret-repl-db/repl-db.secret.sh

# Hack to upgrade database schema
# FIXME improve it
entrypoint --log-level DEBUG smig-update --repl-connection "{{.ReplicationDatabaseURL}}"

entrypoint --log-level DEBUG replication-controller \
    --db-uri "{{.ReplicationDatabaseURL}}" \
    --db-admin-uri "{{.ReplicationDatabaseRootURL}}" \
    --qserv-czar-db="{{.CzarDatabaseRootURL}}" \
    --log-cfg-file "/cm-etc/log.cnf" \
    -- \
    --controller-auto-register-workers=1 \
    --controller-http-server-port="{{.HTTPPort}}" \
    --controller-job-timeout-sec=57600 \
    --controller-request-timeout-sec=57600 \
    --debug \
    --health-probe-interval=120 \
    --instance-id="{{.QservInstance}}" \
    --registry-host="{{.ReplicationRegistryDN}}" \
    --registry-port="{{.HTTPPort}}" \
    --replication-interval=1200 \
    --worker-evict-timeout=3600 \
    --xrootd-host="{{.XrootdRedirectorDN}}" \
    --qserv-sync-force

