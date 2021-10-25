package syncers

import (
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	qservv1beta1 "github.com/lsst/qserv-operator/api/v1beta1"
	"github.com/lsst/qserv-operator/controllers/objects"
	"github.com/lsst/qserv-operator/controllers/syncer"
	"github.com/lsst/qserv-operator/controllers/util"
)

// NewCzarStatefulSetSyncer returns a new sync.Interface for reconciling Qserv Czar StatefulSet
func NewCzarStatefulSetSyncer(r *qservv1beta1.Qserv, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	statefulSet := objects.GenerateCzarStatefulSet(r)
	return syncer.NewObjectSyncer("CzarStatefulSet", r, statefulSet, c, scheme, util.NoFunc)
}

// NewIngestDbStatefulSetSyncer returns a new sync.Interface for reconciling Qserv ingest Db StatefulSet
func NewIngestDbStatefulSetSyncer(r *qservv1beta1.Qserv, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	statefulSet := objects.GenerateIngestDbStatefulSet(r)
	return syncer.NewObjectSyncer("IngestDbStatefulSet", r, statefulSet, c, scheme, util.NoFunc)
}

// NewReplicationCtlStatefulSetSyncer returns a new sync.Interface for reconciling Qserv replication controller StatefulSet
func NewReplicationCtlStatefulSetSyncer(r *qservv1beta1.Qserv, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	statefulSet := objects.GenerateReplicationCtlStatefulSet(r)
	return syncer.NewObjectSyncer("ReplicationCtlStatefulSet", r, statefulSet, c, scheme, util.NoFunc)
}

// NewReplicationDbStatefulSetSyncer returns a new sync.Interface for reconciling Qserv replication Db StatefulSet
func NewReplicationDbStatefulSetSyncer(r *qservv1beta1.Qserv, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	statefulSet := objects.GenerateReplicationDbStatefulSet(r)
	return syncer.NewObjectSyncer("ReplicationDbStatefulSet", r, statefulSet, c, scheme, util.NoFunc)
}

// NewWorkerStatefulSetSyncer returns a new sync.Interface for reconciling Qserv Worker StatefulSet
func NewWorkerStatefulSetSyncer(r *qservv1beta1.Qserv, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	statefulSet := objects.GenerateWorkerStatefulSet(r)
	return syncer.NewObjectSyncer("WorkerStatefulSet", r, statefulSet, c, scheme, util.NoFunc)
}

// NewXrootdStatefulSetSyncer returns a new sync.Interface for reconciling xrootd redirectors cluster StatefulSet
func NewXrootdStatefulSetSyncer(r *qservv1beta1.Qserv, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	statefulSet := objects.GenerateXrootdStatefulSet(r)
	return syncer.NewObjectSyncer("XrootdStatefulSet", r, statefulSet, c, scheme, util.NoFunc)
}
