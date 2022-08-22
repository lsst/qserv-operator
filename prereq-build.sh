#!/bin/bash

# Helper to install operator-sdk 

set -euxo pipefail

RELEASE_VERSION=v1.22.2
export ARCH=$(case $(arch) in x86_64) echo -n amd64 ;; aarch64) echo -n arm64 ;; *) echo -n $(arch) ;; esac)
export OS=$(uname | awk '{print tolower($0)}')

cd /tmp

PGP_SERVER="keyserver.ubuntu.com"
#PGP_SERVER="pool.sks-keyservers.net"
OPERATOR_SDK_DL_URL="https://github.com/operator-framework/operator-sdk/releases/download/${RELEASE_VERSION}"
OPERATOR_SDK_BIN="operator-sdk_${OS}_${ARCH}"

curl -OJL $OPERATOR_SDK_DL_URL/$OPERATOR_SDK_BIN
gpg --keyserver "$PGP_SERVER" --recv-key "052996E2A20B5C7E"
curl -LO ${OPERATOR_SDK_DL_URL}/checksums.txt
curl -LO ${OPERATOR_SDK_DL_URL}/checksums.txt.asc
gpg -u "Operator SDK (release) <cncf-operator-sdk@cncf.io>" --verify checksums.txt.asc
chmod +x "$OPERATOR_SDK_BIN"
sudo mkdir -p /usr/local/bin
sudo cp "$OPERATOR_SDK_BIN" /usr/local/bin/operator-sdk
rm "$OPERATOR_SDK_BIN"
