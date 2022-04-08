#!/bin/bash

# Start Qserv replication registry

# @author  Fabrice Jammes, IN2P3/SLAC

set -exo pipefail

. /secret-mariadb/mariadb.secret.sh
. /secret-repl-db/repl-db.secret.sh

entrypoint --log-level DEBUG replication-registry \
    --db-uri "{{.ReplicationDatabaseURL}}" \
    --db-admin-uri "{{.ReplicationDatabaseRootURL}}" \
    -- \
    --instance-id="{{.QservInstance}}" \
    --registry-port="{{.HTTPPort}}" \
    --debug

