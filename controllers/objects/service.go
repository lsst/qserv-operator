package objects

import (
	qservv1beta1 "github.com/lsst/qserv-operator/api/v1beta1"
	"github.com/lsst/qserv-operator/controllers/constants"
	"github.com/lsst/qserv-operator/controllers/util"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GenerateQservQueryService generate service specification for Qserv Czar proxy
func GenerateQservQueryService(cr *qservv1beta1.Qserv) *v1.Service {
	name := util.GetName(cr, constants.QservName)
	namespace := cr.Namespace

	labels := util.GetComponentLabels(constants.Czar, cr.Name)

	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
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
}

// GenerateCzarService generate service specification for Qserv Czar database
func GenerateCzarService(cr *qservv1beta1.Qserv) *v1.Service {
	name := util.GetCzarServiceName(cr)
	namespace := cr.Namespace

	labels := util.GetComponentLabels(constants.Czar, cr.Name)

	return &v1.Service{
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
				{
					Port:     constants.ProxyPort,
					Protocol: v1.ProtocolTCP,
					Name:     constants.ProxyPortName,
				},
			},
			Selector: labels,
		},
	}
}

// GenerateWorkerService generates headless service for Qserv workers StatefulSet
func GenerateWorkerService(cr *qservv1beta1.Qserv) *v1.Service {
	name := util.GetWorkerServiceName(cr)
	namespace := cr.Namespace

	labels := util.GetComponentLabels(constants.Worker, cr.Name)

	return &v1.Service{
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
}
