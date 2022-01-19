Build qserv-operator
--------------------

Prerequisites
~~~~~~~~~~~~~

-  `git <https://git-scm.com/downloads>`__
-  `go <https://golang.org/dl/>`__ version v1.16+.
-  `docker <https://docs.docker.com/install/>`__ version 20.10+.
-  Access to a Kubernetes v1.20+ cluster. For development purpose, `kind <https://kind.sigs.k8s.io/>`__ is recommended.
   This `kind quickstart <https://github.com/k8s-school/kind-helper#run-kind-on-a-workstation-in-two-lines-of-code>`__ might help.
-  Operator-sdk v1.15+ (see below for quick install)
-  `kubectl <https://kubernetes.io/docs/tasks/tools/install-kubectl/>`__
   version v1.20+.

Build
~~~~~

`qserv-operator` is based on *operator-sdk v1.9.0*

.. code:: sh

    git clone https://github.com/lsst/qserv-operator.git
    cd qserv-operator
    # Pre-requisite, run it only once
    ./install-operator-sdk.sh
    # Build qserv-operator image
    ./build.sh
    # If using `kind`, push qserv-operator image to it
    ./push-image.sh -kd

Test qserv-operator
~~~~~~~~~~~~~~~~~~~

.. code:: sh

    # Install qserv-operator
    ./deploy.sh
    # Install qserv
    kubectl apply -k manifests/base
    ./tests/e2e/integration.sh

Generate and upload documentation
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Documentation is automatically built and generated on each Travis-CI build. This can also be performed manually by launching script below:

.. code:: sh

    curl -fsSL https://raw.githubusercontent.com/lsst-dm/qserv-doc-container/master/run.sh | bash -s -- -p <LTD_PASSWORD> ~/src/qserv
