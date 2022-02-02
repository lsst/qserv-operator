package specs

import (
	qservv1beta1 "github.com/lsst/qserv-operator/api/v1beta1"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ServiceSpec provide default procedures for all Qserv Services specifications
type ServiceSpec struct {
	qserv *qservv1beta1.Qserv
}

// Initialize initialize service specification, default for all Qserv Services specifications
func (c *ServiceSpec) Initialize(qserv *qservv1beta1.Qserv) client.Object {
	c.qserv = qserv
	var object client.Object = &v1.Service{}
	return object
}

// Update update service specification, default for all Qserv Services specifications
func (c *ServiceSpec) Update(object client.Object) (bool, error) {
	return false, nil
}
