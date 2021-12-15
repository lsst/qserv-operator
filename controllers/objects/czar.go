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

type CzarSpec struct {
}

func (c *CzarSpec) Initialize() client.Object {
	var object client.Object = &appsv1.StatefulSet{}
	return object
}

// Create generate statefulset specification for Qserv Czar
func (c *CzarSpec) Create(cr *qservv1beta1.Qserv, object *client.Object) error {
	name := cr.Name + "-" + string(constants.Czar)
	namespace := cr.Namespace
	labels := util.GetComponentLabels(constants.Czar, cr.Name)

	reqLogger := log.WithValues("Request.Namespace", cr.Namespace, "Request.Name", cr.Name)

	storageClass := getValue(cr.Spec.Czar.StorageClass, cr.Spec.StorageClass)
	storageSize := getValue(cr.Spec.Czar.StorageCapacity, cr.Spec.StorageCapacity)

	initContainer, initVolumes := getInitContainer(cr, constants.Czar)
	mariadbContainer, mariadbVolumes := getMariadbContainer(cr, constants.Czar)
	proxyContainer, proxyVolumes := getProxyContainer(cr)

	var volumes VolumeSet
	volumes.make(initVolumes, mariadbVolumes, proxyVolumes)

	ss := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: name,
			Replicas:    &cr.Spec.Czar.Replicas,
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
					Affinity: &cr.Spec.Czar.Affinity,
					InitContainers: []v1.Container{
						initContainer,
					},
					Containers: []v1.Container{
						mariadbContainer,
						proxyContainer,
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
		},
	}

	addDebuggerContainer(reqLogger, ss, cr)

	ss.Spec.Template.Spec.Tolerations = cr.Spec.Tolerations
	*object = ss
	return nil
}

// Update update statefulset specification for Qserv Czar
func (c *CzarSpec) Update(cr *qservv1beta1.Qserv, object *client.Object) (bool, error) {
	return false, nil
}
