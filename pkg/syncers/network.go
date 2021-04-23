package syncers

import (
	qservv1alpha1 "github.com/lsst/qserv-operator/api/v1alpha1"
	qserv "github.com/lsst/qserv-operator/pkg/resources"
	"github.com/lsst/qserv-operator/pkg/syncer"
	"github.com/lsst/qserv-operator/pkg/util"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewNetworkPoliciesSyncer generate NetworkPolicies specifications all Qserv pods
func NewNetworkPoliciesSyncer(r *qservv1alpha1.Qserv, c client.Client, scheme *runtime.Scheme) []syncer.Interface {

	labels := util.GetInstanceLabels(r.Name)
	return []syncer.Interface{
		syncer.NewObjectSyncer("DefaultNetworkPolicy", r, qserv.GenerateDefaultNetworkPolicy(r, labels), c, scheme, util.NoFunc),
		syncer.NewObjectSyncer("CzarNetworkPolicy", r, qserv.GenerateCzarNetworkPolicy(r, labels), c, scheme, util.NoFunc),
		syncer.NewObjectSyncer("ReplDBNetworkPolicy", r, qserv.GenerateReplDBNetworkPolicy(r, labels), c, scheme, util.NoFunc),
		syncer.NewObjectSyncer("WorkerNetworkPolicy", r, qserv.GenerateWorkerNetworkPolicy(r, labels), c, scheme, util.NoFunc),
		syncer.NewObjectSyncer("XrootdRedirectoryNetworkPolicy", r, qserv.GenerateXrootdRedirectorNetworkPolicy(r, labels), c, scheme, util.NoFunc),
	}
}
