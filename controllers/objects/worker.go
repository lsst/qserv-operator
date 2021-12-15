package objects

import (
	qservv1beta1 "github.com/lsst/qserv-operator/api/v1beta1"
	"github.com/lsst/qserv-operator/controllers/constants"
	"github.com/lsst/qserv-operator/controllers/util"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type WorkerSpec struct {
}

func (c *WorkerSpec) GetName() string {
	return string(constants.Worker)
}

func (c *WorkerSpec) Initialize() client.Object {
	var object client.Object = &appsv1.StatefulSet{}
	return object
}

// Create generate statefulset specification for Qserv Czar
func (c *WorkerSpec) Create(cr *qservv1beta1.Qserv, object *client.Object) error {
	name := cr.Name + "-" + c.GetName()
	namespace := cr.Namespace

	labels := util.GetComponentLabels(constants.Worker, cr.Name)

	reqLogger := log.WithValues("Request.Namespace", cr.Namespace, "Request.Name", cr.Name)

	replicas := cr.Spec.Worker.Replicas

	storageClass := getValue(cr.Spec.Worker.StorageClass, cr.Spec.StorageClass)
	storageSize := getValue(cr.Spec.Worker.StorageCapacity, cr.Spec.StorageCapacity)

	initContainer, initVolumes := getInitContainer(cr, constants.Worker)
	mariadbContainer, mariadbVolumes := getMariadbContainer(cr, constants.Worker)
	xrootdContainers, xrootdVolumes := getXrootdContainers(cr, constants.Worker)
	replicationWrkContainer, replicationWrkVolumes := getReplicationWrkContainer(cr)

	// Volumes
	var volumes VolumeSet
	volumes.make(initVolumes, mariadbVolumes, replicationWrkVolumes, xrootdVolumes)

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
					Affinity: &cr.Spec.Worker.Affinity,
					InitContainers: []v1.Container{
						initContainer,
					},
					Containers: []v1.Container{
						mariadbContainer,
						replicationWrkContainer,
						xrootdContainers[0],
						xrootdContainers[1],
					},
					SecurityContext: &v1.PodSecurityContext{
						FSGroup: &constants.QservGID,
					},
					Volumes: volumes.toSlice(),
				},
			},
			VolumeClaimTemplates: []v1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: constants.DataVolumeClaimTemplateName,
					},
					Spec: v1.PersistentVolumeClaimSpec{
						AccessModes:      []v1.PersistentVolumeAccessMode{v1.ReadWriteOnce},
						StorageClassName: &storageClass,
						Resources: v1.ResourceRequirements{
							Requests: v1.ResourceList{
								"storage": resource.MustParse(storageSize),
							},
						},
					},
				},
			},
		}}

	addDebuggerContainer(reqLogger, ss, cr)

	ss.Spec.Template.Spec.Tolerations = cr.Spec.Tolerations

	*object = ss
	return nil
}

// Update update statefulset specification for Qserv Czar
func (c *WorkerSpec) Update(cr *qservv1beta1.Qserv, object *client.Object) (bool, error) {
	return false, nil
}
