package util

import (
	"fmt"

	qservv1alpha1 "github.com/lsst/qserv-operator/pkg/apis/qserv/v1alpha1"
	"github.com/lsst/qserv-operator/pkg/constants"
)

// // GetRedisShutdownConfigMapName returns the name for redis configmap
// func GetRedisShutdownConfigMapName(r *qservv1alpha1.Qserv) string {
//         if r.Spec.Redis.ShutdownConfigMap != "" {
//                 return r.Spec.Redis.ShutdownConfigMap
//         }
//         return GetRedisShutdownName(r)
// }

// // GetRedisName returns the name for redis resources
// func GetRedisName(r *qservv1alpha1.Qserv) string {
//         return generateName(constants.RedisName, r.Name)
// }

// // GetRedisShutdownName returns the name for redis resources
// func GetRedisShutdownName(r *qservv1alpha1.Qserv) string {
//         return generateName(constants.RedisShutdownName, r.Name)
// }

// GetXrootdName returns the name for xrootd resources
func GetXrootdConfigName(r *qservv1alpha1.Qserv) string {
	return generateName(constants.XrootdConfigName, r.Name)
}

func generateName(typeName, metaName string) string {
	return fmt.Sprintf("%s-%s-%s", constants.BaseName, typeName, metaName)
}
