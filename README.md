# qserv-operator

A qserv operator for Kubernetes based on [operator-framework](https://github.com/operator-framework). An Operator is a method of packaging, deploying and managing a Kubernetes application.

## Continuous integration for master branch

Build Qserv-operator and run Qserv multi-node integration tests (using a fixed Qserv version)

| CI       | Status                                                                                                                                                           | Image build  | e2e tests | Documentation generation        | Static code analysis  | Image security scan |
|----------|------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------|-----------|---------------------------------|-----------------------|---------------------|
| Gihub    | [![Qserv CI](https://github.com/lsst/qserv-operator/workflows/CI/badge.svg?branch=master)](https://github.com/lsst/qserv-operator/actions?query=workflow%3A"CI") | Yes          | No        | https://qserv-operator.lsst.io/ | Yes                   | Yes                 |
| Travis   | [![Build Status](https://travis-ci.org/lsst/qserv-operator.svg?branch=master)](https://travis-ci.org/lsst/qserv-operator)                                        | Yes          | Yes (k8s) | No                              | No                    | No                  |

## Documentation

Access to [Qserv-operator documentation](https://qserv-operator.lsst.io/)

# Code analysis

[![Go Report Card](https://goreportcard.com/badge/github.com/xrootd/xrootd-k8s-operator)](https://goreportcard.com/report/github.com/xrootd/xrootd-k8s-operator)

[Security overview](https://github.com/lsst/qserv-operator/security)
