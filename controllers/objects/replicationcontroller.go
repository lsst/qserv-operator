package objects

import (
	qservv1beta1 "github.com/lsst/qserv-operator/api/v1beta1"
	"github.com/lsst/qserv-operator/controllers/constants"
	"github.com/lsst/qserv-operator/controllers/util"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ReplicationControllerSpec struct {
	qserv *qservv1beta1.Qserv
}

func (c *ReplicationControllerSpec) GetName() string {
	return c.qserv.Name + "-" + string(constants.ReplCtl)
}

func (c *ReplicationControllerSpec) Initialize(qserv *qservv1beta1.Qserv) client.Object {
	c.qserv = qserv
	var object client.Object = &appsv1.StatefulSet{}
	return object
}

// Create generate statefulset specification for Qserv Czar
func (c *ReplicationControllerSpec) Create() (client.Object, error) {
	name := c.GetName()
	cr := c.qserv
	namespace := cr.Namespace

	labels := util.GetComponentLabels(constants.ReplCtl, cr.Name)

	reqLogger := log.WithValues("Request.Namespace", cr.Namespace, "Request.Name", cr.Name)

	var replicas int32 = 1

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
			Replicas:            &replicas,
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

// Update update statefulset specification for Qserv Czar
func (c *ReplicationControllerSpec) Update(object client.Object) (bool, error) {
	return false, nil
}

type ReplicationControllerServiceSpec struct {
	qserv *qservv1beta1.Qserv
}

func (c *ReplicationControllerServiceSpec) GetName() string {
	return util.GetReplCtlServiceName(c.qserv)
}

func (c *ReplicationControllerServiceSpec) Initialize(qserv *qservv1beta1.Qserv) client.Object {
	c.qserv = qserv
	var object client.Object = &v1.Service{}
	return object
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
					Port:     constants.ReplicationControllerPort,
					Protocol: v1.ProtocolTCP,
					Name:     constants.ReplicationControllerPortName,
				},
			},
			Selector: labels,
		},
	}
	return service, nil
}

// Update update service specification for Qserv Replication Controller
func (c *ReplicationControllerServiceSpec) Update(object client.Object) (bool, error) {
	return false, nil
}
