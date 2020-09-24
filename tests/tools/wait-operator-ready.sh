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

# Wait Qserv pods to be in ready state

# @author Fabrice Jammes SLAC/IN2P3

set -eux

VERBOSE=false
DIR=$(cd "$(dirname "$0")"; pwd -P)

while test $# -gt 0; do
  case "$1" in
    -v | --verbose)
      VERBOSE=true
      shift
      ;;
  esac
done

while ! kubectl wait deploy -n qserv-operator-system --for=condition=available qserv-operator-controller-manager --timeout="10s"
do
  echo "Wait for Qserv operator to be available"
done

echo "Qserv operator is available"
