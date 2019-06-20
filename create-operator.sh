#!/bin/sh

# See
# https://github.com/operator-framework/operator-sdk/blob/master/doc/user-guide.md#build-and-run-the-operator

set -e
set -x

export GO111MODULE=on
operator-sdk new qserv-operator
cd qserv-operator
operator-sdk add api --api-version=qserv.lsst.org/v1alpha1 --kind=Qserv
operator-sdk add controller --api-version=qserv.lsst.org/v1alpha1 --kind=Qserv
