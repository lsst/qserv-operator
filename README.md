# qserv-operator
A qserv-operator based on operator-sdk
(operator-sdk version operator-sdk version: v0.8.1, commit: 33b3bfe10176f8647f5354516fff29dea42b6342)

## Build

### Prerequisites

- [git][git_tool]
- [go][go_tool] version v1.12+.
- [docker][docker_tool] version 17.03+.
- [kubectl][kubectl_tool] version v1.11.3+.
- Access to a Kubernetes v1.14.2+ cluster.

### Build qserv-operator

```sh
$ git clone https://github.com/kube-incubator/qserv-operator.git
$ cd qserv-operator
$ ./build.sh 
```
## Deploy

### Deploy qserv-operator

```sh
$ ./deploy.sh 
```

### Deploy sample qserv cluster

TODO

### Connect to the qserv cluster

TODO

[git_tool]:https://git-scm.com/downloads
[go_tool]:https://golang.org/dl/
[docker_tool]:https://docs.docker.com/install/
[kubectl_tool]:https://kubernetes.io/docs/tasks/tools/install-kubectl/
