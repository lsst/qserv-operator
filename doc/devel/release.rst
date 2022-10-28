Bump `Qserv` and `Mariadb` versions in `qserv-operator`
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

#. In `qserv` source repository run:

    .. code:: sh

        qserv env | grep -w QSERV_IMAGE

    This command display the Qserv image tag for the current commit in the `qserv` source repository.

    Replace the value of the `spec.image` field image tag in the file `qserv-operator/manifests/base/image.yaml` with the above  Qserv image tag

#. Replace `QSERV_IMAGE` with `QSERV_MARIADB_IMAGE` and `spec.image` with `spec.dbImage` in the above procedure to upagrade the Mariadb container image for this Qserv version.

#. Commit and push to Github the changes for `qserv-operator/manifests/base/image.yaml`  and check if `CI integration tests <https://github.com/lsst/qserv-operator/actions/>`__ are successfull
    If yes, it is possible then to build a release.


Build a Qserv release
~~~~~~~~~~~~~~~~~~~~~

Once CI integration tests are successful for `qserv-operator`.

#. Create an annotated tag for Qserv, named `RELEASE_TAG`, based on the following template `YYYY.M.D-rcX`, and upgrade `Qserv`/`Mariadb` image name in `qserv-operator`, and commit changed to `qserv-operator/manifests/base/image.yaml` file.
#. Then, create a release for `qserv-operator`, on main branch tip:

    .. code:: sh

        ./qserv-operator/publish-release.sh <RELEASE_TAG>

#. Perform the same operation for `qserv-ingest`, on main branch tip:

    .. code:: sh

        ./qserv-operator/publish-release.sh <RELEASE_TAG>

Double-check that CI integration test are passing for both `qserv-ingest` and `qserv-operator`.


TODO
~~~~
- CI has to create/push an image with correct tag (??): Idea use existing images and retag them
- Manage operatorHub

