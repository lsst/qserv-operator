Build qserv-operator
--------------------

Prerequisites
~~~~~~~~~~~~~

-  `git <https://git-scm.com/downloads>`__
-  `go <https://golang.org/dl/>`__ version v1.12+.
-  `docker <https://docs.docker.com/install/>`__ version 17.03+.
-  `kubectl <https://kubernetes.io/docs/tasks/tools/install-kubectl/>`__
   version v1.11.3+.
-  Access to a Kubernetes v1.14.2+ cluster.

Build
~~~~~

`qserv-operator` is based on *operator-sdk v0.15.2*

.. code:: sh

    git clone https://github.com/lsst/qserv-operator.git
    cd qserv-operator
    ./build-all.sh

Test qserv-operator
~~~~~~~~~~~~~~~~~~~

.. code:: sh

    ./deploy/qserv.sh --dev --install-kubedb
    ./run-multinode-tests.sh

Generate and upload documentation
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Documentation is automatically built and generated on each Travis-CI build. This can also be performed manually by launching script below:
    
.. code:: sh

    curl -fsSL https://raw.githubusercontent.com/lsst/doc-container/master/run.sh | bash -s -- -p <LTD_PASSWORD> ~/src/qserv