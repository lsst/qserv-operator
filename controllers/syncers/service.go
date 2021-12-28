package syncers

import (
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	qservv1beta1 "github.com/lsst/qserv-operator/api/v1beta1"
	"github.com/lsst/qserv-operator/controllers/objects"
	"github.com/lsst/qserv-operator/controllers/syncer"
	"github.com/lsst/qserv-operator/controllers/util"
)

// NewQservServicesSyncer returns a new []sync.Interface for reconciling all Qserv services
func NewQservServicesSyncer(r *qservv1beta1.Qserv, c client.Client, scheme *runtime.Scheme) []syncer.Interface {
	syncers := []syncer.Interface{
		syncer.NewObjectSyncer("QservQueryService", r, objects.GenerateQservQueryService(r), c, scheme, util.NoFunc),
		syncer.NewObjectSyncer("CzarService", r, objects.GenerateCzarService(r), c, scheme, util.NoFunc),
	}
	return syncers
}
