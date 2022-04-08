package specs

import (
	qservv1beta1 "github.com/lsst/qserv-operator/api/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// StatefulSetSpec provide default procedures for all Qserv Deployment specifications
type DeploymentSpec struct {
	qserv     *qservv1beta1.Qserv
	hasUpdate bool
}

// Initialize initialize Deployment specification, default for all Qserv Deployment specifications
func (c *DeploymentSpec) Initialize(qserv *qservv1beta1.Qserv) client.Object {
	c.qserv = qserv
	var object client.Object = &appsv1.Deployment{}
	return object
}

// Update update Deployment specification, used as default for all Qserv Deployment specifications
func (c *DeploymentSpec) update(object client.Object, replicas int32) (bool, error) {
	deployment := object.(*appsv1.Deployment)
	c.updateContainersImages(deployment)
	c.updateReplicas(deployment, replicas)
	return c.provideUpdate(), nil
}

func (c *DeploymentSpec) updateContainersImages(deployment *appsv1.Deployment) {
	deploymentContainers := deployment.Spec.Template.Spec.Containers
	hasUpdate := updateContainersImages(c.qserv, deploymentContainers)
	if !c.hasUpdate {
		c.hasUpdate = hasUpdate
	}
	stsInitContainers := deployment.Spec.Template.Spec.InitContainers
	hasUpdate = updateContainersImages(c.qserv, stsInitContainers)
	if !c.hasUpdate {
		c.hasUpdate = hasUpdate
	}
}

func (c *DeploymentSpec) updateReplicas(deployment *appsv1.Deployment, replicas int32) {
	if *deployment.Spec.Replicas != replicas {
		deployment.Spec.Replicas = &replicas
		c.hasUpdate = true
	}
}

func (c *DeploymentSpec) provideUpdate() bool {
	return c.hasUpdate
}
