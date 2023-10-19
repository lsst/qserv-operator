# qserv-operator

A qserv operator for Kubernetes based on [operator-framework](https://github.com/operator-framework). An Operator is a method of packaging, deploying and managing a Kubernetes application.

## Continuous integration for master branch

Build Qserv-operator and run Qserv multi-node integration tests (using a fixed Qserv version)

| CI       | Status                                                                                                                                                           | Image build  | e2e tests | Documentation generation        | Static code analysis  | Image security scan |
|----------|------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------|-----------|---------------------------------|-----------------------|---------------------|
| Gihub    | [![Qserv CI](https://github.com/lsst/qserv-operator/workflows/CI/badge.svg?branch=main)](https://github.com/lsst/qserv-operator/actions?query=workflow%3A"CI") | Yes          | Yes        | https://qserv-operator.lsst.io/ | Yes                   | Yes                 |

## Documentation

Access to [Qserv-operator documentation](https://qserv-operator.lsst.io/)

# Code analysis

[![Go Report Card](https://goreportcard.com/badge/github.com/xrootd/xrootd-k8s-operator)](https://goreportcard.com/report/github.com/xrootd/xrootd-k8s-operator)

[Security overview](https://github.com/lsst/qserv-operator/security)

## How to publish a new release for the whole Qserv stack

### Qserv (and _worker, _master flavor containers), qserv_distrib, qserv_testdata

These are built and published by running the two jenkins jobs  [rebuild-publish-qserv-dev](https://ci.lsst.codes/blue/organizations/jenkins/dax%2Frelease%2Frebuild_publish_qserv-dev/activity) and [build-dev](https://ci.lsst.codes/blue/organizations/jenkins/dax%2Fdocker%2Fbuild-dev/activity), after pushing tags to all the involved repositories. Then release tags must be added to the resulting containers on docker hub.

### qserv-operator

Validate the integration of `qserv-operator` with the release in CI (i.e. GHA), using a dedicated branch

```
cd <project_directory>
# RELEASE format is "<YYYY>.<M>.<i>-rc<j>"
RELEASE="2023.10.1-rc1"
git checkout -b $RELEASE
# Script below edit `qserv` image name in `manifests/image.yaml`, and prepare operatorHub packaging
./publish-release.sh "$RELEASE"
```

Then edit `qserv-ingest` version in `tests/e2e/integration.sh`, to validate the release component altogether.
Once the release CI pass, merge the release branch to `main` branch.

In `main` branch, create the release tag and the image
```
git tag -a "$RELEASE" -m "Version $RELEASE"
git push --follow-tags

./push-image.sh
```

This will automatically push the release tag to the repositories, and push the tagged container images to docker hub.

## How to publish a new release to operatorHub

The above step (i.e. release publishing) must have been completed before doing this one.

```
make bundle
RELEASE="2023.10.1-rc1"
OPERATOR_SRC_DIR="$PWD"
# Clone community-operators and create a branch
gh repo clone https://github.com/lsst/community-operators.git /tmp/community-operators
cd /tmp/community-operators
# Synchronize with upstream repository
git fetch upstream
git checkout main
git merge upstream/main
git push origin main
# Prepare a Pull-Request
git checkout -b "$RELEASE"
cp -r $OPERATOR_SRC_DIR/bundle /tmp/community-operators/operators/qserv-operator/"$RELEASE"
# WARNING: Edit manually 'version' and 'replaceVersion' fields at the end of file qserv-operator.clusterserviceversion.yaml
git add .
git commit --signoff -m "Release $RELEASE for qserv-operator"
git push --set-upstream origin "$RELEASE"
gh repo view --web
# Then make a PR: https://github.com/lsst/community-operators/compare
```
---
**NOTE**

If a CI test fail in PR for [community-operators](https://github.com/k8s-operatorhub/community-operators) official repository, it is possible to run it locally on a workstation using:
```
RELEASE="2023.10.1-rc1"
OPP_PRODUCTION_TYPE=k8s bash <(curl -sL https://raw.githubusercontent.com/redhat-openshift-ecosystem/community-operators-pipeline/ci/latest/ci/scripts/opp.sh) \
kiwi operators/qserv-operator/$RELEASE
```
---
