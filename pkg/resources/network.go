package qserv

import (
	qservv1beta1 "github.com/lsst/qserv-operator/api/v1alpha1"
	"github.com/lsst/qserv-operator/pkg/constants"
	"github.com/lsst/qserv-operator/pkg/util"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// GenerateDefaultNetworkPolicy generate a NetworkPolicy
// which prevents all incoming network connection to all pods in namespace
func GenerateDefaultNetworkPolicy(cr *qservv1beta1.Qserv, labels map[string]string) *v1.NetworkPolicy {
	return &v1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default-deny-ingress",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: v1.NetworkPolicySpec{
			PolicyTypes: []v1.PolicyType{
				v1.PolicyTypeIngress,
			},
			PodSelector: metav1.LabelSelector{},
		},
	}
}

// GenerateCzarNetworkPolicy generate a NetworkPolicy for czar pod
func GenerateCzarNetworkPolicy(cr *qservv1beta1.Qserv, labels map[string]string) *v1.NetworkPolicy {
	return &v1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "allow-czar-ingress",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: v1.NetworkPolicySpec{
			PolicyTypes: []v1.PolicyType{
				v1.PolicyTypeIngress,
			},
			PodSelector: metav1.LabelSelector{
				MatchLabels: util.GetComponentLabels(constants.Czar, cr.Name),
			},
			Ingress: []v1.NetworkPolicyIngressRule{
				{
					// DB port
					Ports: []v1.NetworkPolicyPort{
						{
							Port: &intstr.IntOrString{
								IntVal: constants.MariadbPort,
							},
						},
					},
					From: []v1.NetworkPolicyPeer{
						{
							// Only Replication Controller can access the DB
							PodSelector: &metav1.LabelSelector{
								MatchLabels: util.GetComponentLabels(constants.ReplCtl, cr.Name),
							},
						},
					},
				},
				{
					// Proxy port
					Ports: []v1.NetworkPolicyPort{
						{
							Port: &intstr.IntOrString{
								IntVal: constants.ProxyPort,
							},
						},
					},
				},
			},
		},
	}
}

// GenerateReplDBNetworkPolicy generate a NetworkPolicy for replication database pod
func GenerateReplDBNetworkPolicy(cr *qservv1beta1.Qserv, labels map[string]string) *v1.NetworkPolicy {
	return &v1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "allow-repl-db-ingress",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: v1.NetworkPolicySpec{
			PolicyTypes: []v1.PolicyType{
				v1.PolicyTypeIngress,
			},
			PodSelector: metav1.LabelSelector{
				MatchLabels: util.GetComponentLabels(constants.ReplDb, cr.Name),
			},
			Ingress: []v1.NetworkPolicyIngressRule{
				{
					// DB port
					Ports: []v1.NetworkPolicyPort{
						{
							Port: &intstr.IntOrString{
								IntVal: constants.MariadbPort,
							},
						},
					},
					From: []v1.NetworkPolicyPeer{
						// {
						// 	// Only Replication Controller can access the DB
						// 	PodSelector: &metav1.LabelSelector{
						// 		MatchLabels: util.GetLabels(constants.ReplName, cr.Name),
						// 	},
						// },
					},
				},
			},
		},
	}
}

// GenerateWorkerNetworkPolicy generate a NetworkPolicy for worker pods
func GenerateWorkerNetworkPolicy(cr *qservv1beta1.Qserv, labels map[string]string) *v1.NetworkPolicy {
	return &v1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "allow-worker-ingress",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: v1.NetworkPolicySpec{
			PolicyTypes: []v1.PolicyType{
				v1.PolicyTypeIngress,
			},
			PodSelector: metav1.LabelSelector{
				MatchLabels: util.GetComponentLabels(constants.Worker, cr.Name),
			},
			Ingress: []v1.NetworkPolicyIngressRule{
				{
					// DB port
					Ports: []v1.NetworkPolicyPort{
						{
							Port: &intstr.IntOrString{
								IntVal: constants.MariadbPort,
							},
						},
					},
					From: []v1.NetworkPolicyPeer{
						// {
						// 	// Only Replication Controller can access the DB
						// 	PodSelector: &metav1.LabelSelector{
						// 		MatchLabels: util.GetLabels(constants.ReplName, cr.Name),
						// 	},
						// },
					},
				},
				{
					// Xrootd port
					Ports: []v1.NetworkPolicyPort{
						{
							Port: &intstr.IntOrString{
								IntVal: constants.XrootdPort,
							},
						},
					},
				},
			},
		},
	}
}

// GenerateXrootdRedirectorNetworkPolicy generate a NetworkPolicy for xrootd redirector pods
func GenerateXrootdRedirectorNetworkPolicy(cr *qservv1beta1.Qserv, labels map[string]string) *v1.NetworkPolicy {
	return &v1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "allow-xrootd-redirector-ingress",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: v1.NetworkPolicySpec{
			PolicyTypes: []v1.PolicyType{
				v1.PolicyTypeIngress,
			},
			PodSelector: metav1.LabelSelector{
				MatchLabels: util.GetComponentLabels(constants.XrootdRedirector, cr.Name),
			},
			Ingress: []v1.NetworkPolicyIngressRule{
				{
					// CMSD port
					Ports: []v1.NetworkPolicyPort{
						{
							Port: &intstr.IntOrString{
								IntVal: constants.CmsdPort,
							},
						},
					},
					From: []v1.NetworkPolicyPeer{
						{
							// Only Xrootd workers can access the redirector CMSD
							PodSelector: &metav1.LabelSelector{
								MatchLabels: util.GetComponentLabels(constants.Worker, cr.Name),
							},
						},
					},
				},
				{
					// Xrootd port
					Ports: []v1.NetworkPolicyPort{
						{
							Port: &intstr.IntOrString{
								IntVal: constants.XrootdPort,
							},
						},
					},
					From: []v1.NetworkPolicyPeer{
						{
							// Only CZAR can access the redirector Xrootd port
							PodSelector: &metav1.LabelSelector{
								MatchLabels: util.GetComponentLabels(constants.Czar, cr.Name),
							},
						},
					},
				},
			},
		},
	}
}
