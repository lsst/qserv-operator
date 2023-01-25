#!/bin/bash

# Start Qserv replication worker service inside pod

# @author  Fabrice Jammes, IN2P3/SLAC

# WARN: password are displayed in debug logs
set -euxo pipefail

. /secret-mariadb/mariadb.secret.sh
. /secret-repl-db/repl-db.secret.sh

# Hack to upgrade database schema
# FIXME improve it
# TODO check for local mysql startup
entrypoint --log-level DEBUG smig-update --worker-connection "{{.WorkerDatabaseLocalRootURL}}"

entrypoint worker-repl \
  --db-admin-uri "{{.WorkerDatabaseLocalRootURL}}" \
  --repl-connection "{{.ReplicationDatabaseURL}}" \
  --log-cfg-file "/cm-etc/log.cnf" \
  -- \
  --worker-ingest-num-retries={{.WorkerIngestNumRetries}} \
  --worker-ingest-max-retries={{.WorkerIngestMaxRetries}} \
  --registry-host="{{.ReplicationRegistryDN}}" \
  --registry-port={{.HTTPPort}} \
  --instance-id="{{.QservInstance}}" \
  --debug
