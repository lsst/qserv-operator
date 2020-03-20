#!/bin/sh

set -eux

SRC_DIR="/qserv/data"
DEST_DIR="/qserv/qserv-prod"

ls /qserv/data
exit 0

echo "Copying Qserv data"
sudo -u qserv -- rsync -avz --delete "$SRC_DIR" "$DEST_DIR"

SRC_DIR="/qserv/replication"
if [ -d "$SRC_DIR/mysql" ]; then
    echo "Copying replication data"
    sudo -u qserv -- mkdir -p "$DEST_DIR"
    sudo -u qserv -- rsync -avz --delete "$SRC_DIR" "$DEST_DIR"
fi

