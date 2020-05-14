package sync

import (
	qservv1alpha1 "github.com/lsst/qserv-operator/pkg/apis/qserv/v1alpha1"
	"github.com/lsst/qserv-operator/pkg/constants"
	"github.com/lsst/qserv-operator/pkg/scheme/qserv"
	"github.com/lsst/qserv-operator/pkg/staging/syncer"
	"github.com/lsst/qserv-operator/pkg/util"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewNetworkPoliciesSyncer(r *qservv1alpha1.Qserv, c client.Client, scheme *runtime.Scheme) []syncer.Interface {
	labels := util.MergeLabels(controllerLabels, util.GetLabels(constants.NetworkPolicy, r.Name))
	return []syncer.Interface{
		syncer.NewObjectSyncer("DefaultNetworkPolicy", r, qserv.GenerateDefaultNetworkPolicy(r, labels), c, scheme, noFunc),
		syncer.NewObjectSyncer("CzarNetworkPolicy", r, qserv.GenerateCzarNetworkPolicy(r, labels), c, scheme, noFunc),
		syncer.NewObjectSyncer("ReplDBNetworkPolicy", r, qserv.GenerateReplDBNetworkPolicy(r, labels), c, scheme, noFunc),
		syncer.NewObjectSyncer("WorkerNetworkPolicy", r, qserv.GenerateWorkerNetworkPolicy(r, labels), c, scheme, noFunc),
		syncer.NewObjectSyncer("XrootdRedirectoryNetworkPolicy", r, qserv.GenerateXrootdRedirectorNetworkPolicy(r, labels), c, scheme, noFunc),
	}
}
