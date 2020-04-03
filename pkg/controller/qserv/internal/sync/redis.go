package sync

import (
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	qservv1alpha1 "github.com/lsst/qserv-operator/pkg/apis/qserv/v1alpha1"
	"github.com/lsst/qserv-operator/pkg/scheme/qserv"
	"github.com/lsst/qserv-operator/pkg/staging/syncer"
)

// NewRedisSyncer returns a new sync.Interface for reconciling Qserv Redis database, build on top of KubeDb Redis CRD
func NewRedisSyncer(r *qservv1alpha1.Qserv, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	redis := qserv.GenerateRedis(r, controllerLabels)
	return syncer.NewObjectSyncer("Redis", r, redis, c, scheme, noFunc)
}
