package sync

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	qservv1alpha1 "github.com/lsst/qserv-operator/pkg/apis/qserv/v1alpha1"
	"github.com/lsst/qserv-operator/pkg/scheme/qserv"
	"github.com/lsst/qserv-operator/pkg/staging/syncer"
)

func NewMicroserviceConfigMapSyncer(r *qservv1alpha1.Qserv, c client.Client, scheme *runtime.Scheme, microservice string, subpath string) syncer.Interface {
	cm := qserv.GenerateMicroserviceConfigMap(r, controllerLabels, microservice, subpath)
	objectName := fmt.Sprintf("%s%sConfigMap", strings.Title(microservice), strings.Title(subpath))
	return syncer.NewObjectSyncer(objectName, r, cm, c, scheme, func(existing runtime.Object) error {
		return nil
	})
}

func NewDomainNameConfigMapSyncer(r *qservv1alpha1.Qserv, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	cm := qserv.GenerateDomainNameConfigMap(r, controllerLabels)
	return syncer.NewObjectSyncer("DomainNameConfigMap", r, cm, c, scheme, func(existing runtime.Object) error {
		return nil
	})
}

func NewSqlConfigMapSyncer(r *qservv1alpha1.Qserv, c client.Client, scheme *runtime.Scheme, db string) syncer.Interface {
	cm := qserv.GenerateSqlConfigMap(r, controllerLabels, db)
	objectName := fmt.Sprintf("%sSqlConfigMap", strings.Title(db))
	return syncer.NewObjectSyncer(objectName, r, cm, c, scheme, func(existing runtime.Object) error {
		return nil
	})
}

// // NewRedisShutdownConfigMapSyncer returns a new sync.Interface for reconciling Redis Shutdown ConfigMap
// func NewRedisShutdownConfigMapSyncer(r *redisv1alpha1.Redis, c client.Client, scheme *runtime.Scheme) syncer.Interface {
// 	cm := redis.GenerateRedisShutdownConfigMap(r, controllerLabels)
// 	return syncer.NewObjectSyncer("RedisShutdownConfigMap", r, cm, c, scheme, func(existing runtime.Object) error {
// 		return nil
// 	})
// }

// // NewSentinelConfigMapSyncer returns a new sync.Interface for reconciling Sentinel ConfigMap
// func NewSentinelConfigMapSyncer(r *redisv1alpha1.Redis, c client.Client, scheme *runtime.Scheme) syncer.Interface {
// 	cm := redis.GenerateSentinelConfigMap(r, controllerLabels)
// 	return syncer.NewObjectSyncer("SentinelConfigMap", r, cm, c, scheme, func(existing runtime.Object) error {
// 		return nil
// 	})
// }
