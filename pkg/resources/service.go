package qserv

import (
	qservv1alpha1 "github.com/lsst/qserv-operator/api/v1alpha1"
	"github.com/lsst/qserv-operator/pkg/constants"
	"github.com/lsst/qserv-operator/pkg/util"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GenerateQservNodePortService generate NodePort service specification for Qserv Czar proxy
func GenerateQservNodePortService(cr *qservv1alpha1.Qserv, labels map[string]string) *v1.Service {
	name := util.GetName(cr, constants.QservName)
	namespace := cr.Namespace

	labels = util.MergeLabels(labels, util.GetLabels(constants.Czar, cr.Name))

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

// GenerateCzarService generate service specification for Qserv Czar database
func GenerateCzarService(cr *qservv1alpha1.Qserv, labels map[string]string) *v1.Service {
	name := util.GetCzarServiceName(cr)
	namespace := cr.Namespace

	labels = util.MergeLabels(labels, util.GetLabels(constants.Czar, cr.Name))

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

// GenerateIngestDbService generate service specification for Qserv Ingest database
func GenerateIngestDbService(cr *qservv1alpha1.Qserv, labels map[string]string) *v1.Service {
	name := util.GetName(cr, string(constants.IngestDb))
	namespace := cr.Namespace

	labels = util.MergeLabels(labels, util.GetLabels(constants.IngestDb, cr.Name))

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

// GenerateReplicationCtlService generate service specification for Qserv Replication Controller
func GenerateReplicationCtlService(cr *qservv1alpha1.Qserv, labels map[string]string) *v1.Service {
	name := util.GetReplCtlServiceName(cr)
	namespace := cr.Namespace

	labels = util.MergeLabels(labels, util.GetLabels(constants.ReplCtl, cr.Name))

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
					Port:     constants.ReplicationControllerPort,
					Protocol: v1.ProtocolTCP,
					Name:     constants.ReplicationControllerPortName,
				},
			},
			Selector: labels,
		},
	}
}

// GenerateReplicationDbService generate service specification for Qserv Replication Controller database
func GenerateReplicationDbService(cr *qservv1alpha1.Qserv, labels map[string]string) *v1.Service {
	name := util.GetName(cr, string(constants.ReplDb))
	namespace := cr.Namespace

	labels = util.MergeLabels(labels, util.GetLabels(constants.ReplDb, cr.Name))

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

	labels = util.MergeLabels(labels, util.GetLabels(constants.Worker, cr.Name))

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
	name := util.GetName(cr, string(constants.XrootdRedirector))
	namespace := cr.Namespace

	labels = util.MergeLabels(labels, util.GetLabels(constants.XrootdRedirector, cr.Name))

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
