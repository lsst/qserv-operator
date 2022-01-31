package specs

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
	qserv           *qservv1beta1.Qserv
	qservContainers []constants.ContainerName
}

func (c *WorkerSpec) GetName() string {
	return util.GetName(c.qserv, string(constants.Worker))
}

func (c *WorkerSpec) Initialize(qserv *qservv1beta1.Qserv) client.Object {
	c.qserv = qserv
	c.qservContainers = []constants.ContainerName{constants.CmsdName, constants.XrootdName, constants.ReplWrkName}
	var object client.Object = &appsv1.StatefulSet{}
	return object
}

// Create generate statefulset specification for Qserv Czar
func (c *WorkerSpec) Create() (client.Object, error) {
	name := c.GetName()
	cr := c.qserv
	namespace := cr.Namespace

	labels := util.GetComponentLabels(constants.Worker, cr.Name)

	reqLogger := log.WithValues("Request.Namespace", cr.Namespace, "Request.Name", cr.Name)

	replicas := cr.Spec.Worker.Replicas

	storageClass := util.GetValue(cr.Spec.Worker.StorageClass, cr.Spec.StorageClass)
	storageSize := util.GetValue(cr.Spec.Worker.StorageCapacity, cr.Spec.StorageCapacity)

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

	return ss, nil
}

// Update update statefulset specification for Qserv Worker
func (c *WorkerSpec) Update(object client.Object) (bool, error) {
	image := c.qserv.Spec.Worker.Image
	ss := object.(*appsv1.StatefulSet)

	ssContainers := ss.Spec.Template.Spec.Containers
	hasUpdate := updateContainersImages(image, ssContainers, c.qservContainers)

	// TODO add support for the below feature, which is
	// currently forbidden by admissionWebhook
	// Ensure the deployment size is the same as the spec.
	replicas := c.qserv.Spec.Worker.Replicas

	if *ss.Spec.Replicas != replicas {
		ss.Spec.Replicas = &replicas
		hasUpdate = true
	}

	return hasUpdate, nil
}

// WorkerServiceSpec allows to reconcile xrootd service
type WorkerServiceSpec struct {
	qserv *qservv1beta1.Qserv
}

func (c *WorkerServiceSpec) GetName() string {
	return util.GetName(c.qserv, string(constants.Worker))
}

func (c *WorkerServiceSpec) Initialize(qserv *qservv1beta1.Qserv) client.Object {
	c.qserv = qserv
	var object client.Object = &v1.Service{}
	return object
}

// Create generates headless service for Qserv workers StatefulSet
func (c *WorkerServiceSpec) Create() (client.Object, error) {
	cr := c.qserv
	name := c.GetName()
	namespace := cr.Namespace

	labels := util.GetComponentLabels(constants.Worker, cr.Name)

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
					Port:     constants.XrootdPort,
					Protocol: v1.ProtocolTCP,
					Name:     constants.XrootdPortName,
				},
			},
			Selector: labels,
		},
	}
	return service, nil
}

// Update update service specification for Qserv workers
func (c *WorkerServiceSpec) Update(object client.Object) (bool, error) {
	return false, nil
}
