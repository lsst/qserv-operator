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

# Wait Qserv statefulset to be in running state

# @author Fabrice Jammes SLAC/IN2P3

set -e
set -x

DIR=$(cd "$(dirname "$0")"; pwd -P)

OUT="/tmp/out.txt"
STR="INFO:lsst.qserv.admin.commons:stdout file: '/qserv/run/tmp/qservTest_case01/outputs/qserv/0001.1_fetchObjectById.txt'"
"$DIR/run-integration-tests.sh" >& "$OUT" &

echo "Wait for error to happen"

until grep "$STR" "$OUT" 
do
    echo -n "Wait for error"
    sleep 2
echo -n '.'
done

echo "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
kubectl logs example-qserv-czar-0 -c proxy
echo "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
kubectl logs example-qserv-xrootd-redirector-0 -c xrootd
echo "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
kubectl logs example-qserv-xrootd-redirector-0 -c cmsd
echo "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
kubectl logs example-qserv-worker-0 -c xrootd
echo "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
kubectl logs example-qserv-worker-0 -c cmsd
echo "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
kubectl logs example-qserv-worker-1 -c xrootd
echo "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
kubectl logs example-qserv-worker-1 -c cmsd

tail -f "$OUT"
