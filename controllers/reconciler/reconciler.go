package reconciler

import (
	qservv1beta1 "github.com/lsst/qserv-operator/api/v1beta1"
	"github.com/lsst/qserv-operator/controllers/objects"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ReconcileFunc func(qserv *qservv1beta1.Qserv, object *client.Object) error

type ObjectSpec struct {
	Create ReconcileFunc
	Update ReconcileFunc
}

func nilFunc(qserv *qservv1beta1.Qserv, object *client.Object) error {
	return nil
}

// NewCzarStatefulSetSyncer returns a new sync.Interface for reconciling Qserv Czar StatefulSet
func NewCzarStatefulSetSyncer() *ObjectSpec {
	return &ObjectSpec{
		Create: objects.GenerateCzarStatefulSet,
		Update: nilFunc,
	}
}
