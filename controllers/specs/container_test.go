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

func TestGetInitContainer(t *testing.T) {

	qserv := &qservv1beta1.Qserv{}
	qserv.Name = "qserv"
	qserv.Spec.Image = "qserv/qserv:nil"
	qserv.Spec.StorageCapacity = "10Gi"
	qserv.Spec.Worker.Replicas = 1
	pod := constants.Czar
	container, volumes := getInitContainer(qserv, pod)
	// containerByte, _ := yaml.Marshal(&container)
	// t.Logf("Container: %s", containerByte)
	// volumeByte, _ := yaml.Marshal(&volumes)
	// t.Logf("Volume: %s", volumeByte)

	t.Logf("VolumeMounts: %v", container.VolumeMounts)

	expectedVolumeMount := [...]string{"qserv-data", "config-mariadb-etc", "config-initdb-start", "secret-mariadb", "config-sql-czar"}
	volumeMountNames := getNames(container.VolumeMounts)
	for _, n := range expectedVolumeMount {
		assert.Contains(t, volumeMountNames, n)
	}

	t.Logf("Volumes: %v", volumes.getNames())
	expectedVolumeNames := [...]string{"config-sql-czar", "config-mariadb-etc", "config-initdb-start", "secret-mariadb-qserv"}
	volumeNames := volumes.getNames()
	for _, n := range expectedVolumeNames {
		assert.Contains(t, volumeNames, n)
	}

}

func TestGetMariadbContainer(t *testing.T) {

	qserv := &qservv1beta1.Qserv{}
	qserv.Name = "qserv"
	qserv.Spec.Image = "qserv/qserv:nil"
	qserv.Spec.StorageCapacity = "10Gi"
	qserv.Spec.Worker.Replicas = 1
	pod := constants.Czar
	container, volumes := getMariadbContainer(qserv, pod)
	// containerByte, _ := yaml.Marshal(&container)
	// t.Logf("Container: %s", containerByte)
	// volumeByte, _ := yaml.Marshal(&volumes)
	// t.Logf("Volume: %s", volumeByte)

	t.Logf("VolumeMounts: %v", container.VolumeMounts)

	expectedVolumeNames := [...]string{"config-mariadb-etc", "config-mariadb-start", "tmp-volume"}
	for _, n := range expectedVolumeNames {
		assert.Contains(t, volumes.getNames(), n)
	}

}

func TestGetReplicationRegistryContainer(t *testing.T) {

	qserv := &qservv1beta1.Qserv{}
	qserv.Name = "qserv"
	qserv.Spec.Image = "qserv/qserv:nil"
	qserv.Spec.StorageCapacity = "10Gi"
	qserv.Spec.Worker.Replicas = 1
	container, volumes := getReplicationRegistryContainer(qserv)
	// containerByte, _ := yaml.Marshal(&container)
	// t.Logf("Container: %s", containerByte)
	// volumeByte, _ := yaml.Marshal(&volumes)
	// t.Logf("Volume: %s", volumeByte)

	t.Logf("VolumeMounts: %v", container.VolumeMounts)

	expectedVolumeNames := [...]string{"config-repl-registry-start", "config-repl-registry-etc", "secret-repl-db-qserv", "secret-mariadb-qserv"}
	for _, n := range expectedVolumeNames {
		assert.Contains(t, volumes.getNames(), n)
	}

}
