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
}

func (c *ReplicationControllerSpec) GetName() string {
	return string(constants.ReplCtl)
}

func (c *ReplicationControllerSpec) Initialize() client.Object {
	var object client.Object = &appsv1.StatefulSet{}
	return object
}

// Create generate statefulset specification for Qserv Czar
func (c *ReplicationControllerSpec) Create(cr *qservv1beta1.Qserv) (client.Object, error) {
	name := cr.Name + "-" + c.GetName()
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
func (c *ReplicationControllerSpec) Update(cr *qservv1beta1.Qserv, object client.Object) (bool, error) {
	return false, nil
}
