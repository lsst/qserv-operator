package specs

import (
	"github.com/lsst/qserv-operator/controllers/constants"
	"github.com/lsst/qserv-operator/controllers/util"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// IngestDatabaseSpec provide procedures for Ingest Database StatefulSet specification
type IngestDatabaseSpec struct {
	StatefulSetSpec
}

// GetName return name for Ingest Database StatefulSet
func (c *IngestDatabaseSpec) GetName() string {
	return c.qserv.Name + "-" + string(constants.IngestDb)
}

// Create generate statefulset specification for Qserv Ingest Database
func (c *IngestDatabaseSpec) Create() (client.Object, error) {
	cr := c.qserv
	name := c.GetName()
	namespace := cr.Namespace

	labels := util.GetComponentLabels(constants.IngestDb, cr.Name)

	storageClass := cr.Spec.StorageClass
	storageSize := cr.Spec.StorageCapacity

	initContainer, initVolumes := getInitContainer(cr, constants.IngestDb)
	mariadbContainer, mariadbVolumes := getMariadbContainer(cr, constants.IngestDb)

	var volumes VolumeSet
	volumes.make(initVolumes, mariadbVolumes)

	ss := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName: name,
			Replicas:    &constants.IngestDatabaseReplicas,
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
					Affinity: &cr.Spec.Ingest.Affinity,
					InitContainers: []v1.Container{
						initContainer,
					},
					Containers: []v1.Container{
						mariadbContainer,
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
						StorageClassName: storageClass,
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

	ss.Spec.Template.Spec.Tolerations = cr.Spec.Tolerations

	return ss, nil
}

// Update update  Ingest Database specification
func (c *IngestDatabaseSpec) Update(object client.Object) (bool, error) {
	return c.update(object, constants.IngestDatabaseReplicas)
}

// IngestDatabaseServiceSpec allows to reconcile Ingest Database Service
type IngestDatabaseServiceSpec struct {
	ServiceSpec
}

// GetName return name for Ingest Database Service
func (c *IngestDatabaseServiceSpec) GetName() string {
	return util.GetName(c.qserv, string(constants.IngestDb))
}

// Create generate service specification for Qserv Ingest database
func (c *IngestDatabaseServiceSpec) Create() (client.Object, error) {
	cr := c.qserv
	name := c.GetName()
	namespace := cr.Namespace

	labels := util.GetComponentLabels(constants.IngestDb, cr.Name)

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
					Port:     constants.MariadbPort,
					Protocol: v1.ProtocolTCP,
					Name:     constants.MariadbPortName,
				},
			},
			Selector: labels,
		},
	}
	return service, nil
}
