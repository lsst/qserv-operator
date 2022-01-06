package specs

import (
	"fmt"
	"testing"

	qservv1beta1 "github.com/lsst/qserv-operator/api/v1beta1"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
)

func TestUpdate(t *testing.T) {

	spec := WorkerSpec{}
	qserv := &qservv1beta1.Qserv{}

	qserv.Name = "qserv"
	qserv.Spec.StorageCapacity = "10Gi"
	qserv.Spec.Worker.Replicas = 1

	spec.Initialize(qserv)
	object, _ := spec.Create()

	qserv.Spec.Worker.Replicas = 1
	spec.Update(object)

	ss := object.(*appsv1.StatefulSet)

	fmt.Printf("ss %v\n", ss)

	assert.Equal(t, qserv.Spec.Worker.Replicas, *((*ss).Spec.Replicas))
}
