package syncers

import (
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	qservv1beta1 "github.com/lsst/qserv-operator/api/v1beta1"
	"github.com/lsst/qserv-operator/controllers/objects"
	"github.com/lsst/qserv-operator/controllers/syncer"
	"github.com/lsst/qserv-operator/controllers/util"
)

// NewReplicationDbStatefulSetSyncer returns a new sync.Interface for reconciling Qserv replication Db StatefulSet
func NewReplicationDbStatefulSetSyncer(r *qservv1beta1.Qserv, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	statefulSet := objects.GenerateReplicationDbStatefulSet(r)
	return syncer.NewObjectSyncer("ReplicationDbStatefulSet", r, statefulSet, c, scheme, util.NoFunc)
}

// NewXrootdStatefulSetSyncer returns a new sync.Interface for reconciling xrootd redirectors cluster StatefulSet
func NewXrootdStatefulSetSyncer(r *qservv1beta1.Qserv, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	statefulSet := objects.GenerateXrootdStatefulSet(r)
	return syncer.NewObjectSyncer("XrootdStatefulSet", r, statefulSet, c, scheme, util.NoFunc)
}
