# Parameters related to NCSA Qserv cluster

# All worker hosts have same prefix
HOSTNAME_TPL="lsst-qserv-db"

INSTANCE='qserv-dev'

# First and last id for worker node names
WORKER_FIRST_ID=01
WORKER_LAST_ID=30

MASTER="lsst-qserv-master01"

WORKERS=$(seq --format "${HOSTNAME_TPL}%02g" \
    --separator=" " "$WORKER_FIRST_ID" "$WORKER_LAST_ID")

PARALLEL_SSH_CFG="$HOME/.ssh/sshloginfile"
