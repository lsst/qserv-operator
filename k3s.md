## Deploy qserv on top of k3s

### Quick start for Ubuntu LTS

```
sudo apt-get update
sudo apt-get install curl docker.io git vim
# then add current user to docker group and restart gnome session
sudo usermod -a -G docker $(id -nu)

curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="--docker --write-kubeconfig-mode 644" sh -
export KUBECONFIG=/etc/rancher/k3s/k3s.yaml

curl -fsSL https://raw.githubusercontent.com/lsst/qserv-operator/tickets\/DM-24372/deploy/qserv.sh | bash -s
kubectl apply -k overlays/ci-redis-k3s 
./wait-qserv-ready.sh
./run-integration-tests.sh
```

