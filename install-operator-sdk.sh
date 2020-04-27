#!/bin/sh

# Helper to install operator-sdk 

set -e
set -x

RELEASE_VERSION=v0.18.1

curl -OJL https://github.com/operator-framework/operator-sdk/releases/download/${RELEASE_VERSION}/operator-sdk-${RELEASE_VERSION}-x86_64-linux-gnu
curl -OJL https://github.com/operator-framework/operator-sdk/releases/download/${RELEASE_VERSION}/operator-sdk-${RELEASE_VERSION}-x86_64-linux-gnu.asc
gpg --keyserver keyserver.ubuntu.com --recv-key "BF6F6F18846753754CBB1DDFBC9679ED89ED8983"
gpg --verify operator-sdk-${RELEASE_VERSION}-x86_64-linux-gnu.asc
chmod +x operator-sdk-${RELEASE_VERSION}-x86_64-linux-gnu
sudo mkdir -p /usr/local/bin
sudo cp operator-sdk-${RELEASE_VERSION}-x86_64-linux-gnu /usr/local/bin/operator-sdk
rm operator-sdk-${RELEASE_VERSION}-x86_64-linux-gnu
