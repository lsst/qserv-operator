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

// ReplicationDatabaseSpec provide procedures for Replication Database StatefulSet specification
type ReplicationDatabaseSpec struct {
	StatefulSetSpec
}

// GetName return name for Replication Database StatefulSet
func (c *ReplicationDatabaseSpec) GetName() string {
	return util.GetName(c.qserv, string(constants.ReplDbName))
}

// Create generate statefulset specification for Qserv Replication Database
func (c *ReplicationDatabaseSpec) Create() (client.Object, error) {
	cr := c.qserv
	name := c.GetName()
	namespace := cr.Namespace

	labels := util.GetComponentLabels(constants.ReplDb, cr.Name)

	storageClass := cr.Spec.StorageClass
	storageSize := cr.Spec.StorageCapacity

	initContainer, initVolumes := getInitContainer(cr, constants.ReplDb)
	mariadbContainer, mariadbVolumes := getMariadbContainer(cr, constants.ReplDb)

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
			Replicas:    &constants.ReplicationDatabaseReplicas,
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

	ss.Spec.Template.Spec.Tolerations = cr.Spec.Tolerations

	return ss, nil
}

// Update update Replication Database specification
func (c *ReplicationDatabaseSpec) Update(object client.Object) (bool, error) {
	return c.update(object, constants.ReplicationDatabaseReplicas)
}

// ReplicationDatabaseServiceSpec allows to reconcile Replication Database Service
type ReplicationDatabaseServiceSpec struct {
	ServiceSpec
}

// GetName return name for Replication Database Service
func (c *ReplicationDatabaseServiceSpec) GetName() string {
	return util.GetName(c.qserv, string(constants.ReplDb))
}

// Create generate service specification for Qserv Replication Controller database
func (c *ReplicationDatabaseServiceSpec) Create() (client.Object, error) {
	cr := c.qserv

	labels := util.GetComponentLabels(constants.ReplDb, cr.Name)

	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.GetName(),
			Namespace: cr.Namespace,
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
