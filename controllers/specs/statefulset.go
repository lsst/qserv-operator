package specs

import (
	qservv1beta1 "github.com/lsst/qserv-operator/api/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// StatefulSetSpec provide default procedures for all Qserv StatefulSetSpec specifications
type StatefulSetSpec struct {
	qserv     *qservv1beta1.Qserv
	hasUpdate bool
}

// Initialize initialize StatefulSet specification, default for all Qserv StatefulSet specifications
func (c *StatefulSetSpec) Initialize(qserv *qservv1beta1.Qserv) client.Object {
	c.qserv = qserv
	var object client.Object = &appsv1.StatefulSet{}
	return object
}

// Update update StatefulSet specification, used as default for all Qserv StatefulSet specifications
func (c *StatefulSetSpec) update(object client.Object, replicas int32) (bool, error) {
	sts := object.(*appsv1.StatefulSet)
	c.updateContainersImages(sts)
	c.updateReplicas(sts, replicas)
	return c.provideUpdate(), nil
}

func (c *StatefulSetSpec) updateContainersImages(sts *appsv1.StatefulSet) {
	stsContainers := sts.Spec.Template.Spec.Containers
	hasUpdate := updateContainersImages(c.qserv, stsContainers)
	if !c.hasUpdate {
		c.hasUpdate = hasUpdate
	}
	stsInitContainers := sts.Spec.Template.Spec.InitContainers
	hasUpdate = updateContainersImages(c.qserv, stsInitContainers)
	if !c.hasUpdate {
		c.hasUpdate = hasUpdate
	}
}

func (c *StatefulSetSpec) updateReplicas(sts *appsv1.StatefulSet, replicas int32) {
	if *sts.Spec.Replicas != replicas {
		sts.Spec.Replicas = &replicas
		c.hasUpdate = true
	}
}

func (c *StatefulSetSpec) provideUpdate() bool {
	return c.hasUpdate
}
