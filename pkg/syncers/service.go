package syncers

import (
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	qservv1alpha1 "github.com/lsst/qserv-operator/api/v1alpha1"
	qserv "github.com/lsst/qserv-operator/pkg/resources"
	"github.com/lsst/qserv-operator/pkg/syncer"
	"github.com/lsst/qserv-operator/pkg/util"
)

// NewDashboardServiceSyncer returns a new sync.Interface for reconciling Ingest Database Service
func NewDashboardServiceSyncer(r *qservv1alpha1.Qserv, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	svc := qserv.GenerateDashboardService(r)
	return syncer.NewObjectSyncer("DashboardService", r, svc, c, scheme, util.NoFunc)
}

// NewIngestDbServiceSyncer returns a new sync.Interface for reconciling Ingest Database Service
func NewIngestDbServiceSyncer(r *qservv1alpha1.Qserv, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	svc := qserv.GenerateIngestDbService(r)
	return syncer.NewObjectSyncer("IngestDbService", r, svc, c, scheme, util.NoFunc)
}

// NewQservServicesSyncer returns a new []sync.Interface for reconciling all Qserv services
func NewQservServicesSyncer(r *qservv1alpha1.Qserv, c client.Client, scheme *runtime.Scheme) []syncer.Interface {
	syncers := []syncer.Interface{
		syncer.NewObjectSyncer("QservQueryService", r, qserv.GenerateQservQueryService(r), c, scheme, util.NoFunc),
		syncer.NewObjectSyncer("CzarService", r, qserv.GenerateCzarService(r), c, scheme, util.NoFunc),
		syncer.NewObjectSyncer("WorkerService", r, qserv.GenerateWorkerService(r), c, scheme, util.NoFunc),
	}
	return syncers
}

// NewReplicationCtlServiceSyncer returns a new sync.Interface for reconciling Replication Controller Service
func NewReplicationCtlServiceSyncer(r *qservv1alpha1.Qserv, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	svc := qserv.GenerateReplicationCtlService(r)
	return syncer.NewObjectSyncer("ReplicationCtlService", r, svc, c, scheme, util.NoFunc)
}

// NewReplicationDbServiceSyncer returns a new sync.Interface for reconciling Replication Database Service
func NewReplicationDbServiceSyncer(r *qservv1alpha1.Qserv, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	svc := qserv.GenerateReplicationDbService(r)
	return syncer.NewObjectSyncer("ReplicationDbService", r, svc, c, scheme, util.NoFunc)
}

// NewXrootdRedirectorServiceSyncer returns a new sync.Interface for reconciling Xrootd Redirector Service
func NewXrootdRedirectorServiceSyncer(r *qservv1alpha1.Qserv, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	svc := qserv.GenerateXrootdRedirectorService(r)
	return syncer.NewObjectSyncer("XrootdRedirectorService", r, svc, c, scheme, util.NoFunc)
}
