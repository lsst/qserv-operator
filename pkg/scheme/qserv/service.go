package qserv

import (
	qservv1alpha1 "github.com/lsst/qserv-operator/pkg/apis/qserv/v1alpha1"
	"github.com/lsst/qserv-operator/pkg/constants"
	"github.com/lsst/qserv-operator/pkg/util"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// func GenerateSentinelService(r *qservv1alpha1.Qserv, labels map[string]string) *corev1.Service {
// 	name := util.GetSentinelName(r)
// 	namespace := r.Namespace

// 	sentinelTargetPort := intstr.FromInt(26379)
// 	labels = util.MergeLabels(labels, util.GetLabels(constants.SentinelRoleName, r.Name))

// 	return &corev1.Service{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      name,
// 			Namespace: namespace,
// 			Labels:    labels,
// 		},
// 		Spec: corev1.ServiceSpec{
// 			Selector: labels,
// 			Ports: []corev1.ServicePort{
// 				{
// 					Name:       "sentinel",
// 					Port:       26379,
// 					TargetPort: sentinelTargetPort,
// 					Protocol:   "TCP",
// 				},
// 			},
// 		},
// 	}
// }

func GenerateWorkerService(cr *qservv1alpha1.Qserv, labels map[string]string) *v1.Service {
	name := util.GetWorkerName(cr)
	namespace := cr.Namespace

	labels = util.MergeLabels(labels, util.GetLabels(constants.WorkerName, cr.Name))

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
					Port:     constants.WmgrPort,
					Protocol: v1.ProtocolTCP,
					Name:     constants.WmgrPortName,
				},
			},
			Selector: labels,
		},
	}
}

func GenerateCzarService(cr *qservv1alpha1.Qserv, labels map[string]string) *v1.Service {
	name := util.GetCzarName(cr)
	namespace := cr.Namespace

	labels = util.MergeLabels(labels, util.GetLabels(constants.CzarName, cr.Name))

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
					Port:     constants.CmsdPort,
					Protocol: v1.ProtocolTCP,
					Name:     constants.CmsdPortName,
				},
			},
			Selector: labels,
		},
	}
}

func GenerateXrootdRedirectorService(cr *qservv1alpha1.Qserv, labels map[string]string) *v1.Service {
	name := util.GetXrootdRedirectorName(cr)
	namespace := cr.Namespace

	labels = util.MergeLabels(labels, util.GetLabels(constants.XrootdRedirectorName, cr.Name))

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
