package specs

import (
	qservv1beta1 "github.com/lsst/qserv-operator/api/v1beta1"
	"github.com/lsst/qserv-operator/controllers/constants"
	"github.com/lsst/qserv-operator/controllers/util"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type XrootdSpec struct {
	qserv *qservv1beta1.Qserv
}

func (c *XrootdSpec) GetName() string {
	return util.GetName(c.qserv, string(constants.XrootdRedirector))
}

func (c *XrootdSpec) Initialize(qserv *qservv1beta1.Qserv) client.Object {
	c.qserv = qserv
	var object client.Object = &appsv1.StatefulSet{}
	return object
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

// Update update statefulset specification for Qserv Ingest Database
func (c *XrootdSpec) Update(object client.Object) (bool, error) {
	return false, nil
}

// XrootdServiceSpec allows to reconcile xrootd service
type XrootdServiceSpec struct {
	qserv *qservv1beta1.Qserv
}

func (c *XrootdServiceSpec) GetName() string {
	return util.GetName(c.qserv, string(constants.XrootdRedirector))
}

func (c *XrootdServiceSpec) Initialize(qserv *qservv1beta1.Qserv) client.Object {
	c.qserv = qserv
	var object client.Object = &v1.Service{}
	return object
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

// Update update service specification for Qserv Replication Controller
func (c *XrootdServiceSpec) Update(object client.Object) (bool, error) {
	ss := object.(*appsv1.StatefulSet)

	ssContainers := ss.Spec.Template.Spec.Containers
	hasUpdate := updateContainersImages(c.qserv, ssContainers)
	return hasUpdate, nil
}
