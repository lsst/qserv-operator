#!/bin/bash

# Run static code analysis for Qserv

# @author  Fabrice Jammes, IN2P3

set -euxo pipefail

DIR=$(cd "$(dirname "$0")"; pwd -P)

BINPATH="$(go env GOPATH)/bin"

go get -u \
	github.com/kisielk/errcheck \
	golang.org/x/tools/cmd/goimports \
	golang.org/x/lint/golint \
	github.com/securego/gosec/cmd/gosec \
	golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow \
	honnef.co/go/tools/cmd/staticcheck

golint -set_exit_status "$DIR"/...
go vet -vettool="$BINPATH/shadow"
staticcheck "$DIR"/...
