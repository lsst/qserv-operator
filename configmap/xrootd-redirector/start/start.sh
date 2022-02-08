#!/bin/sh

# Start cmsd and xrootd inside pod
# Launch as qserv user

# @author  Fabrice Jammes, IN2P3

set -eux

entrypoint xrootd-manager \
          --cmsd-manager-name "{{.XrootdRedirectorDN}}" \
          --cmsd-manager-count "{{.XrootdRedirectorReplicas}}"