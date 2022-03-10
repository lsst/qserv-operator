#!/bin/sh

# Start cmsd redirector inside container
# Launch as qserv user

# @author  Fabrice Jammes, IN2P3

set -eux

entrypoint cmsd-manager --cms-delay-servers {{.WorkerReplicas}}
