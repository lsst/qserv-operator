package qserv

import (
	"fmt"
	qservv1alpha1 "github.com/lsst/qserv-operator/pkg/apis/qserv/v1alpha1"
	"github.com/lsst/qserv-operator/pkg/constants"
	"github.com/lsst/qserv-operator/pkg/util"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GenerateCzarProxyService generate service specification for Qserv Czar proxy
func GenerateCzarProxyService(cr *qservv1alpha1.Qserv, labels map[string]string) *v1.Service {
	suffix := fmt.Sprintf("%s-%s", constants.CzarName, constants.ProxyName)
	name := util.GetName(cr, suffix)
	namespace := cr.Namespace

	labels = util.MergeLabels(labels, util.GetLabels(constants.CzarName, cr.Name))

	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeNodePort,
			Ports: []v1.ServicePort{
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

// GenerateCzarDatabaseService generate service specification for Qserv Czar database
func GenerateCzarDatabaseService(cr *qservv1alpha1.Qserv, labels map[string]string) *v1.Service {
	suffix := fmt.Sprintf("%s-%s", constants.CzarName, constants.MariadbName)
	name := util.GetName(cr, suffix)
	namespace := cr.Namespace

	labels = util.MergeLabels(labels, util.GetLabels(constants.CzarName, cr.Name))

	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeClusterIP,
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
}

// GenerateReplicationCtlService generate service specification for Qserv Czar database
func GenerateReplicationCtlService(cr *qservv1alpha1.Qserv, labels map[string]string) *v1.Service {
	name := util.GetReplCtlServiceName(cr)
	namespace := cr.Namespace

	labels = util.MergeLabels(labels, util.GetLabels(constants.ReplName, cr.Name))

	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeNodePort,
			Ports: []v1.ServicePort{
				{
					Port:     constants.ReplicationControllerPort,
					Protocol: v1.ProtocolTCP,
					Name:     constants.ReplicationControllerPortName,
				},
			},
			Selector: labels,
		},
	}
}

func GenerateReplicationDbService(cr *qservv1alpha1.Qserv, labels map[string]string) *v1.Service {
	name := util.GetName(cr, string(constants.ReplDbName))
	namespace := cr.Namespace

	labels = util.MergeLabels(labels, util.GetLabels(constants.ReplName, cr.Name))

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
			},
			Selector: labels,
		},
	}
}

// GenerateWorkerService generates headless service for Qserv workers StatefulSet
func GenerateWorkerService(cr *qservv1alpha1.Qserv, labels map[string]string) *v1.Service {
	name := util.GetWorkerServiceName(cr)
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

func GenerateXrootdRedirectorService(cr *qservv1alpha1.Qserv, labels map[string]string) *v1.Service {
	name := util.GetName(cr, string(constants.XrootdRedirectorName))
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
