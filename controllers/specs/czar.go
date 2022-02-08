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

// CzarSpec provide procedures for Czar StatefulSet specification
type CzarSpec struct {
	StatefulSetSpec
}

// GetName return name for Czar StatefulSet
func (c *CzarSpec) GetName() string {
	return util.GetName(c.qserv, string(constants.Czar))
}

// Create generate statefulset specification for Qserv Czar
func (c *CzarSpec) Create() (client.Object, error) {
	name := c.GetName()
	cr := c.qserv
	namespace := cr.Namespace
	labels := util.GetComponentLabels(constants.Czar, cr.Name)

	log.Info("Create czar specification")

	storageClass := util.GetValue(cr.Spec.Czar.StorageClass, cr.Spec.StorageClass)
	storageSize := util.GetValue(cr.Spec.Czar.StorageCapacity, cr.Spec.StorageCapacity)

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

	addDebuggerContainer(log, ss, cr)

	ss.Spec.Template.Spec.Tolerations = cr.Spec.Tolerations
	return ss, nil
}

// Update update  Ingest Database specification
func (c *CzarSpec) Update(object client.Object) (bool, error) {
	replicas := c.qserv.Spec.Czar.Replicas
	return c.update(object, replicas)
}

// CzarServiceSpec provide procedures for Czar Service specification
type CzarServiceSpec struct {
	ServiceSpec
}

// GetName return name for Czar Service
func (c *CzarServiceSpec) GetName() string {
	return util.GetName(c.qserv, string(constants.Czar))
}

// Create generate service specification for Qserv Czar
func (c *CzarServiceSpec) Create() (client.Object, error) {
	cr := c.qserv

	labels := util.GetComponentLabels(constants.Czar, cr.Name)

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
				{
					Port:     constants.ProxyPort,
					Protocol: v1.ProtocolTCP,
					Name:     constants.ProxyPortName,
				},
			},
			Selector: labels,
		},
	}
	return service, nil
}

// QueryServiceSpec provide procedures for Query Service specification
type QueryServiceSpec struct {
	ServiceSpec
}

// GetName return name for Query Service
func (c *QueryServiceSpec) GetName() string {
	return util.GetName(c.qserv, string(constants.QservName))
}

// Create generate service specification for Qserv Czar proxy
func (c *QueryServiceSpec) Create() (client.Object, error) {
	cr := c.qserv

	labels := util.GetComponentLabels(constants.Czar, cr.Name)

	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        c.GetName(),
			Namespace:   cr.Namespace,
			Labels:      labels,
			Annotations: cr.Spec.QueryService.Annotations,
		},
		Spec: v1.ServiceSpec{
			LoadBalancerIP: cr.Spec.QueryService.LoadBalancerIP,
			Type:           cr.Spec.QueryService.ServiceType,
			Ports: []v1.ServicePort{
				{
					Name:     constants.ProxyPortName,
					NodePort: cr.Spec.QueryService.NodePort,
					Port:     constants.ProxyPort,
					Protocol: v1.ProtocolTCP,
				},
			},
			Selector: labels,
		},
	}
	return service, nil
}
