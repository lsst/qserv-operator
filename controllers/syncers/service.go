package syncers

import (
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	qservv1beta1 "github.com/lsst/qserv-operator/api/v1beta1"
	"github.com/lsst/qserv-operator/controllers/objects"
	"github.com/lsst/qserv-operator/controllers/syncer"
	"github.com/lsst/qserv-operator/controllers/util"
)

// NewIngestDbServiceSyncer returns a new sync.Interface for reconciling Ingest Database Service
func NewIngestDbServiceSyncer(r *qservv1beta1.Qserv, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	svc := objects.GenerateIngestDbService(r)
	return syncer.NewObjectSyncer("IngestDbService", r, svc, c, scheme, util.NoFunc)
}

// NewQservServicesSyncer returns a new []sync.Interface for reconciling all Qserv services
func NewQservServicesSyncer(r *qservv1beta1.Qserv, c client.Client, scheme *runtime.Scheme) []syncer.Interface {
	syncers := []syncer.Interface{
		syncer.NewObjectSyncer("QservQueryService", r, objects.GenerateQservQueryService(r), c, scheme, util.NoFunc),
		syncer.NewObjectSyncer("CzarService", r, objects.GenerateCzarService(r), c, scheme, util.NoFunc),
		syncer.NewObjectSyncer("WorkerService", r, objects.GenerateWorkerService(r), c, scheme, util.NoFunc),
	}
	return syncers
}

// NewReplicationDbServiceSyncer returns a new sync.Interface for reconciling Replication Database Service
func NewReplicationDbServiceSyncer(r *qservv1beta1.Qserv, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	svc := objects.GenerateReplicationDbService(r)
	return syncer.NewObjectSyncer("ReplicationDbService", r, svc, c, scheme, util.NoFunc)
}

// NewXrootdRedirectorServiceSyncer returns a new sync.Interface for reconciling Xrootd Redirector Service
func NewXrootdRedirectorServiceSyncer(r *qservv1beta1.Qserv, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	svc := objects.GenerateXrootdRedirectorService(r)
	return syncer.NewObjectSyncer("XrootdRedirectorService", r, svc, c, scheme, util.NoFunc)
}
