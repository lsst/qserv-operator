package specs

import (
	"github.com/lsst/qserv-operator/controllers/constants"
	"github.com/lsst/qserv-operator/controllers/util"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ReplicationRegistrySpec provide procedures for Replication Registry Deployment specification
// See https://confluence.lsstcorp.org/display/DM/Configuring+worker+registry+in+the+Replication+system+of+Qserv
type ReplicationRegistrySpec struct {
	DeploymentSpec
}

// GetName return name for Replication Registry Deployment
func (c *ReplicationRegistrySpec) GetName() string {
	return util.GetName(c.qserv, string(constants.ReplRegistry))
}

// Create generate Deployment specification for Replication Registry
func (c *ReplicationRegistrySpec) Create() (client.Object, error) {
	name := c.GetName()
	cr := c.qserv
	namespace := cr.Namespace

	labels := util.GetComponentLabels(constants.ReplRegistry, cr.Name)

	// reqLogger := log.WithValues("Request.Namespace", cr.Namespace, "Request.Name", cr.Name)

	replRegistryContainer, replRegistryVolumes := getReplicationRegistryContainer(cr)

	var volumes VolumeSet
	volumes.make(replRegistryVolumes)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &constants.ReplicationRegistryReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: v1.PodSpec{
					Affinity: &cr.Spec.Replication.Affinity,
					Containers: []v1.Container{
						replRegistryContainer,
					},
					Volumes: volumes.toSlice(),
				},
			},
		},
	}

	deployment.Spec.Template.Spec.Tolerations = cr.Spec.Tolerations

	return deployment, nil
}

// Update update replication registry specification
func (c *ReplicationRegistrySpec) Update(object client.Object) (bool, error) {
	return c.update(object, constants.ReplicationControllerReplicas)
}

// ReplicationRegistryServiceSpec provide procedures for Replication Registry Service specification
type ReplicationRegistryServiceSpec struct {
	ServiceSpec
}

// GetName return name for Replication Registry Service
func (c *ReplicationRegistryServiceSpec) GetName() string {
	return util.GetReplicationRegistryServiceName(c.qserv)
}

// Create generate service specification for Qserv Registry Controller
func (c *ReplicationRegistryServiceSpec) Create() (client.Object, error) {
	name := c.GetName()
	cr := c.qserv
	namespace := cr.Namespace

	labels := util.GetComponentLabels(constants.ReplRegistry, cr.Name)

	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: v1.ServiceSpec{
			Type:      v1.ServiceTypeClusterIP,
			ClusterIP: v1.ClusterIPNone,
			Ports: []v1.ServicePort{
				{
					Port:     constants.HTTPPort,
					Protocol: v1.ProtocolTCP,
					Name:     constants.HTTPPortName,
				},
			},
			Selector: labels,
		},
	}
	return service, nil
}
