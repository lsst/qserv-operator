package syncers

import (
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	qservv1alpha1 "github.com/lsst/qserv-operator/api/v1alpha1"
	qserv "github.com/lsst/qserv-operator/pkg/resources"
	"github.com/lsst/qserv-operator/pkg/syncer"
)

// NewIngestDbServiceSyncer returns a new sync.Interface for reconciling Ingest Database Service
func NewIngestDbServiceSyncer(r *qservv1alpha1.Qserv, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	svc := qserv.GenerateIngestDbService(r, controllerLabels)
	return syncer.NewObjectSyncer("IngestDbService", r, svc, c, scheme, noFunc)
}

// NewQservServicesSyncer returns a new []sync.Interface for reconciling all Qserv services
func NewQservServicesSyncer(r *qservv1alpha1.Qserv, c client.Client, scheme *runtime.Scheme) []syncer.Interface {
	syncers := []syncer.Interface{
		syncer.NewObjectSyncer("QservNodePortService", r, qserv.GenerateQservQueryService(r, controllerLabels), c, scheme, noFunc),
		syncer.NewObjectSyncer("Czar", r, qserv.GenerateCzarService(r, controllerLabels), c, scheme, noFunc),
		syncer.NewObjectSyncer("WorkerService", r, qserv.GenerateWorkerService(r, controllerLabels), c, scheme, noFunc),
	}

	return syncers
}

// NewReplicationCtlServiceSyncer returns a new sync.Interface for reconciling Replication Controller Service
func NewReplicationCtlServiceSyncer(r *qservv1alpha1.Qserv, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	svc := qserv.GenerateReplicationCtlService(r, controllerLabels)
	return syncer.NewObjectSyncer("ReplicationCtlService", r, svc, c, scheme, noFunc)
}

// NewReplicationDbServiceSyncer returns a new sync.Interface for reconciling Replication Database Service
func NewReplicationDbServiceSyncer(r *qservv1alpha1.Qserv, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	svc := qserv.GenerateReplicationDbService(r, controllerLabels)
	return syncer.NewObjectSyncer("ReplicationDbService", r, svc, c, scheme, noFunc)
}

// NewXrootdRedirectorServiceSyncer returns a new sync.Interface for reconciling Xrootd Redirector Service
func NewXrootdRedirectorServiceSyncer(r *qservv1alpha1.Qserv, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	svc := qserv.GenerateXrootdRedirectorService(r, controllerLabels)
	return syncer.NewObjectSyncer("XrootdRedirectorService", r, svc, c, scheme, noFunc)
}
