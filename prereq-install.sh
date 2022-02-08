#!/bin/bash

# Install pre-requisites for deploying Qserv

# @author Fabrice Jammes IN2P3

set -euxo pipefail

VERSION="v1.6.1"
echo "Install cert-manager $VERSION"
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/$VERSION/cert-manager.yaml

VERSION="4.0.5"
echo "Install kustomize $VERSION"
curl -lO "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"
chmod +x ./install_kustomize.sh
sudo rm -f /usr/local/bin/kustomize
sudo ./install_kustomize.sh "$VERSION" /usr/local/bin