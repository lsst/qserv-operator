package specs

import (
	"github.com/lsst/qserv-operator/controllers/constants"
	"github.com/lsst/qserv-operator/controllers/util"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ReplicationControllerSpec provide procedures for Replication Controller StatefulSet specification
type ReplicationControllerSpec struct {
	StatefulSetSpec
}

// GetName return name for Replication Controller StatefulSet
func (c *ReplicationControllerSpec) GetName() string {
	return util.GetName(c.qserv, string(constants.ReplCtl))
}

// Create generate Statefulset specification for Replication Controller
func (c *ReplicationControllerSpec) Create() (client.Object, error) {
	name := c.GetName()
	cr := c.qserv
	namespace := cr.Namespace

	labels := util.GetComponentLabels(constants.ReplCtl, cr.Name)

	reqLogger := log.WithValues("Request.Namespace", cr.Namespace, "Request.Name", cr.Name)

	replCtlContainer, replCtlVolumes := getReplicationCtlContainer(cr)

	var volumes VolumeSet
	volumes.make(replCtlVolumes)

	ss := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			PodManagementPolicy: "Parallel",
			ServiceName:         name,
			Replicas:            &constants.ReplicationControllerReplicas,
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
				Type: "RollingUpdate",
			},
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
						replCtlContainer,
					},
					Volumes: volumes.toSlice(),
				},
			},
		},
	}

	addDebuggerContainer(reqLogger, ss, cr)

	ss.Spec.Template.Spec.Tolerations = cr.Spec.Tolerations

	return ss, nil
}

// Update update replication controller specification
func (c *ReplicationControllerSpec) Update(object client.Object) (bool, error) {
	return c.update(object, constants.ReplicationControllerReplicas)
}

// ReplicationControllerServiceSpec provide procedures for Replication Controller Service specification
type ReplicationControllerServiceSpec struct {
	ServiceSpec
}

// GetName return name for Replication Controller Service
func (c *ReplicationControllerServiceSpec) GetName() string {
	return util.GetReplCtlServiceName(c.qserv)
}

// Create generate service specification for Qserv Replication Controller
func (c *ReplicationControllerServiceSpec) Create() (client.Object, error) {
	name := c.GetName()
	cr := c.qserv
	namespace := cr.Namespace

	labels := util.GetComponentLabels(constants.ReplCtl, cr.Name)

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
