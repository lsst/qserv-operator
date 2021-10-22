package syncers

import (
	qservv1beta1 "github.com/lsst/qserv-operator/api/v1beta1"
	"github.com/lsst/qserv-operator/controllers/objects"
	"github.com/lsst/qserv-operator/controllers/syncer"
	"github.com/lsst/qserv-operator/controllers/util"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewNetworkPoliciesSyncer generate NetworkPolicies specifications all Qserv pods
func NewNetworkPoliciesSyncer(r *qservv1beta1.Qserv, c client.Client, scheme *runtime.Scheme) []syncer.Interface {

	labels := util.GetInstanceLabels(r.Name)
	return []syncer.Interface{
		syncer.NewObjectSyncer("DefaultNetworkPolicy", r, objects.GenerateDefaultNetworkPolicy(r, labels), c, scheme, util.NoFunc),
		syncer.NewObjectSyncer("CzarNetworkPolicy", r, objects.GenerateCzarNetworkPolicy(r, labels), c, scheme, util.NoFunc),
		syncer.NewObjectSyncer("ReplDBNetworkPolicy", r, objects.GenerateReplDBNetworkPolicy(r, labels), c, scheme, util.NoFunc),
		syncer.NewObjectSyncer("WorkerNetworkPolicy", r, objects.GenerateWorkerNetworkPolicy(r, labels), c, scheme, util.NoFunc),
		syncer.NewObjectSyncer("XrootdRedirectoryNetworkPolicy", r, objects.GenerateXrootdRedirectorNetworkPolicy(r, labels), c, scheme, util.NoFunc),
	}
}
