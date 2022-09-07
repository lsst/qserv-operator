#!/bin/bash

# LSST Data Management System
# Copyright 2014 LSST Corporation.
#
# This product includes software developed by the
# LSST Project (http://www.lsst.org/).
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the LSST License Statement and
# the GNU General Public License along with this program.  If not,
# see <http://www.lsstcorp.org/LegalNotices/>.
# Launch Qserv multinode tests on Swarm cluster

# @author Fabrice Jammes SLAC/IN2P3

set -euxo pipefail

INGEST_DIR="/tmp/qserv-ingest"
INGEST_VERSION="main"
INSTANCE=$(kubectl get qservs.qserv.lsst.org -o=jsonpath='{.items[0].metadata.name}')

echo "Run integration tests for Qserv"
git clone https://github.com/lsst-dm/qserv-ingest.git  --branch "$INGEST_VERSION" --single-branch "$INGEST_DIR"
git -C "$INGEST_DIR" checkout -b ci
"$INGEST_DIR"/prereq-install.sh
"$INGEST_DIR"/itest/run-tests.sh
