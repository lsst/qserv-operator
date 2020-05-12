package qserv

import (
	qservv1alpha1 "github.com/lsst/qserv-operator/pkg/apis/qserv/v1alpha1"
	"github.com/lsst/qserv-operator/pkg/constants"
	"github.com/lsst/qserv-operator/pkg/util"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GenerateDefaultNetworkPolicy(cr *qservv1alpha1.Qserv, labels map[string]string) *v1.NetworkPolicy {
	namespace := cr.Namespace
	labels = util.MergeLabels(labels, util.GetLabels(constants.NetworkPolicy, cr.Name))

	return &v1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default-deny-ingress",
			Namespace: namespace,
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
