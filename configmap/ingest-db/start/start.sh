#!/bin/bash

# Start mariadb inside pod
# and do not exit

# @author  Fabrice Jammes, IN2P3/SLAC

set -eux

echo "-- Start mariadb server."
mysqld
if [ $? -ne 0 ]; then
    >&2 echo "ERROR: failed to start the replication database"
    exit 1
fi
