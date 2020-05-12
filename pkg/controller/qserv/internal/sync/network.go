package sync

import (
	qservv1alpha1 "github.com/lsst/qserv-operator/pkg/apis/qserv/v1alpha1"
	"github.com/lsst/qserv-operator/pkg/scheme/qserv"
	"github.com/lsst/qserv-operator/pkg/staging/syncer"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewNetworkPoliciesSyncer(r *qservv1alpha1.Qserv, c client.Client, scheme *runtime.Scheme) []syncer.Interface {
	policy := qserv.GenerateDefaultNetworkPolicy(r, controllerLabels)
	return []syncer.Interface{
		syncer.NewObjectSyncer("DefaultNetworkPolicy", r, policy, c, scheme, noFunc),
	}
}
