package specs

import (
	qservv1beta1 "github.com/lsst/qserv-operator/api/v1beta1"
	"github.com/lsst/qserv-operator/controllers/constants"
	"github.com/lsst/qserv-operator/controllers/util"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NetworkPolicy which prevents all incoming network connection to all pods in namespace
type DefaultNetworkPolicySpec struct {
	qserv *qservv1beta1.Qserv
}

func (c *DefaultNetworkPolicySpec) GetName() string {
	return "default-deny-ingress"
}

func (c *DefaultNetworkPolicySpec) Initialize(qserv *qservv1beta1.Qserv) client.Object {
	c.qserv = qserv
	var object client.Object = &v1.NetworkPolicy{}
	return object
}

func (c *DefaultNetworkPolicySpec) Create() (client.Object, error) {
	cr := c.qserv
	labels := util.GetInstanceLabels(cr.Name)
	networkPolicy := &v1.NetworkPolicy{
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
	return networkPolicy, nil
}

func (c *DefaultNetworkPolicySpec) Update(object client.Object) (bool, error) {
	return false, nil
}

// NetworkPolicy for Czar Pod
type CzarNetworkPolicySpec struct {
	qserv *qservv1beta1.Qserv
}

func (c *CzarNetworkPolicySpec) GetName() string {
	return "allow-czar-ingress"
}

func (c *CzarNetworkPolicySpec) Initialize(qserv *qservv1beta1.Qserv) client.Object {
	c.qserv = qserv
	var object client.Object = &v1.NetworkPolicy{}
	return object
}

func (c *CzarNetworkPolicySpec) Create() (client.Object, error) {
	cr := c.qserv
	labels := util.GetInstanceLabels(cr.Name)
	networkPolicy := &v1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.GetName(),
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
	return networkPolicy, nil
}

func (c *CzarNetworkPolicySpec) Update(object client.Object) (bool, error) {
	return false, nil
}

// NetworkPolicy for Replication Database Pod
type ReplDatabaseNetworkPolicySpec struct {
	qserv *qservv1beta1.Qserv
}

func (c *ReplDatabaseNetworkPolicySpec) GetName() string {
	return "allow-repl-db-ingress"
}

func (c *ReplDatabaseNetworkPolicySpec) Initialize(qserv *qservv1beta1.Qserv) client.Object {
	c.qserv = qserv
	var object client.Object = &v1.NetworkPolicy{}
	return object
}

func (c *ReplDatabaseNetworkPolicySpec) Create() (client.Object, error) {
	cr := c.qserv
	labels := util.GetInstanceLabels(cr.Name)
	networkPolicy := &v1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.GetName(),
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
	return networkPolicy, nil
}

func (c *ReplDatabaseNetworkPolicySpec) Update(object client.Object) (bool, error) {
	return false, nil
}

// NetworkPolicy for Worker Pods
type WorkerNetworkPolicySpec struct {
	qserv *qservv1beta1.Qserv
}

func (c *WorkerNetworkPolicySpec) GetName() string {
	return "allow-worker-ingress"
}

func (c *WorkerNetworkPolicySpec) Initialize(qserv *qservv1beta1.Qserv) client.Object {
	c.qserv = qserv
	var object client.Object = &v1.NetworkPolicy{}
	return object
}

func (c *WorkerNetworkPolicySpec) Create() (client.Object, error) {
	cr := c.qserv
	labels := util.GetInstanceLabels(cr.Name)
	networkPolicy := &v1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.GetName(),
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
	return networkPolicy, nil
}

func (c *WorkerNetworkPolicySpec) Update(object client.Object) (bool, error) {
	return false, nil
}

// NetworkPolicy for xrootd redirector pods
type XrootdRedirectorNetworkPolicySpec struct {
	qserv *qservv1beta1.Qserv
}

func (c *XrootdRedirectorNetworkPolicySpec) GetName() string {
	return "allow-xrootd-redirector-ingress"
}

func (c *XrootdRedirectorNetworkPolicySpec) Initialize(qserv *qservv1beta1.Qserv) client.Object {
	c.qserv = qserv
	var object client.Object = &v1.NetworkPolicy{}
	return object
}

func (c *XrootdRedirectorNetworkPolicySpec) Create() (client.Object, error) {
	cr := c.qserv
	labels := util.GetInstanceLabels(cr.Name)
	networkPolicy := &v1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.GetName(),
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
	return networkPolicy, nil
}

func (c *XrootdRedirectorNetworkPolicySpec) Update(object client.Object) (bool, error) {
	return false, nil
}
