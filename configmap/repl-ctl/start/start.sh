#!/bin/bash

# Start Qserv replication controller service inside pod

# @author  Fabrice Jammes, IN2P3/SLAC

set -exo pipefail

# Load parameters of the setup into the corresponding environment
# variables

# Start master controller

. /secret-mariadb/mariadb.secret.sh
. /secret-repl-db/repl-db.secret.sh

# Hack to upgrade database schema
# FIXME improve it
entrypoint --log-level DEBUG smig-update --repl-connection "{{.ReplicationDatabaseURL}}"

entrypoint --log-level DEBUG replication-controller \
    --db-uri "{{.ReplicationDatabaseURL}}" \
    --db-admin-uri "{{.ReplicationDatabaseRootURL}}" \
{{- range $val := Iterate .WorkerReplicas}}{{$workerFQDN := print $.WorkerDN "-" $val "." $.WorkerDN}}
    --worker qserv_worker_db=mysql://qsmaster@{{$workerFQDN}}:3306/qservw_worker,host={{$workerFQDN}} \
{{- end}}
    --qserv-czar-db="{{.CzarDatabaseRootURL}}" \
    -- \
    --instance-id="{{.QservInstance}}" \
    --xrootd-host="{{.XrootdRedirectorDN}}" \
    --controller-http-server-port="{{.ReplicationControllerPort}}" \
    --controller-request-timeout-sec=57600 \
    --controller-job-timeout-sec=57600 \
    --worker-evict-timeout=3600 \
    --health-probe-interval=120 \
    --replication-interval=1200

