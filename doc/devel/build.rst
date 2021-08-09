Build qserv-operator
--------------------

Prerequisites
~~~~~~~~~~~~~

-  `git <https://git-scm.com/downloads>`__
-  `go <https://golang.org/dl/>`__ version v1.16+.
-  `docker <https://docs.docker.com/install/>`__ version 20.10+.
-  `kubectl <https://kubernetes.io/docs/tasks/tools/install-kubectl/>`__
   version v1.20+.
-  Access to a Kubernetes v1.20+ cluster.

Build
~~~~~

`qserv-operator` is based on *operator-sdk v1.9.0*

.. code:: sh

    git clone https://github.com/lsst/qserv-operator.git
    cd qserv-operator
    ./build.sh

Test qserv-operator
~~~~~~~~~~~~~~~~~~~

.. code:: sh

    ./deploy.sh
    kubectl apply -k manifests/base
    ./tests/e2e/integration.sh

Generate and upload documentation
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Documentation is automatically built and generated on each Travis-CI build. This can also be performed manually by launching script below:

.. code:: sh

    curl -fsSL https://raw.githubusercontent.com/lsst-dm/qserv-doc-container/master/run.sh | bash -s -- -p <LTD_PASSWORD> ~/src/qserv
