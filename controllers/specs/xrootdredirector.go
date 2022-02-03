package specs

import (
	"github.com/lsst/qserv-operator/controllers/constants"
	"github.com/lsst/qserv-operator/controllers/util"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// XrootdSpec provide procedures for Xrootd Redirector StatefulSet specification
type XrootdSpec struct {
	StatefulSetSpec
}

// GetName return name for Xrootd Redirector StatefulSet
func (c *XrootdSpec) GetName() string {
	return util.GetName(c.qserv, string(constants.XrootdRedirector))
}

// Create generate statefulset specification for xrootd redirectors
func (c *XrootdSpec) Create() (client.Object, error) {
	cr := c.qserv
	namespace := cr.Namespace
	name := c.GetName()

	reqLogger := log.WithValues("Request.Namespace", namespace, "Request.Name", cr.Name)

	labels := util.GetComponentLabels(constants.XrootdRedirector, cr.Name)

	var replicas int32 = cr.Spec.Xrootd.Replicas

	containers, volumes := getXrootdContainers(cr, constants.XrootdRedirector)

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
					Affinity:   &cr.Spec.Xrootd.Affinity,
					Containers: containers,
					Volumes:    volumes.toSlice(),
				},
			},
		},
	}

	addDebuggerContainer(reqLogger, ss, cr)

	ss.Spec.Template.Spec.Tolerations = cr.Spec.Tolerations

	return ss, nil
}

// Update update Xrootd specification
func (c *XrootdSpec) Update(object client.Object) (bool, error) {
	replicas := c.qserv.Spec.Xrootd.Replicas
	return c.update(object, replicas)
}

// XrootdServiceSpec allows to reconcile xrootd service
type XrootdServiceSpec struct {
	ServiceSpec
}

// GetName return name for Xrootd Redirector Service
func (c *XrootdServiceSpec) GetName() string {
	return util.GetName(c.qserv, string(constants.XrootdRedirector))
}

// Create generates headless service specification for xrootd redirectors StatefulSet
func (c *XrootdServiceSpec) Create() (client.Object, error) {
	cr := c.qserv
	name := c.GetName()
	namespace := cr.Namespace

	labels := util.GetComponentLabels(constants.XrootdRedirector, cr.Name)

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
				{
					Port:     constants.CmsdPort,
					Protocol: v1.ProtocolTCP,
					Name:     constants.CmsdPortName,
				},
			},
			Selector: labels,
		},
	}
	return service, nil
}
