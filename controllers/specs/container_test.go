package specs

import (
	"testing"

	qservv1beta1 "github.com/lsst/qserv-operator/api/v1beta1"
	"github.com/lsst/qserv-operator/controllers/constants"
	v1 "k8s.io/api/core/v1"

	"github.com/stretchr/testify/assert"
)

func TestUpdateContainersImages(t *testing.T) {

	qserv := &qservv1beta1.Qserv{}
	qserv.Name = "qserv"
	qserv.Spec.Image = "qserv/qserv:nil"
	qserv.Spec.StorageCapacity = "10Gi"
	qserv.Spec.Worker.Replicas = 1
	initContainer, _ := getInitContainer(qserv, constants.Worker)
	mariadbContainer, _ := getMariadbContainer(qserv, constants.Worker)
	xrootdContainers, _ := getXrootdContainers(qserv, constants.Worker)
	replicationWrkContainer, _ := getReplicationWrkContainer(qserv)

	containers := []v1.Container{
		initContainer,
		mariadbContainer,
		replicationWrkContainer,
		xrootdContainers[0],
		xrootdContainers[1],
	}

	updateContainersImages(qserv, containers)

	assert.NotEqual(t, qserv.Spec.Image, containers[0].Image)
	assert.NotEqual(t, qserv.Spec.Image, containers[1].Image)
	assert.Equal(t, qserv.Spec.Image, containers[2].Image)
	assert.Equal(t, qserv.Spec.Image, containers[3].Image)
	assert.Equal(t, qserv.Spec.Image, containers[4].Image)
}
