package util

import (
	"fmt"

	qservv1alpha1 "github.com/lsst/qserv-operator/pkg/apis/qserv/v1alpha1"
	"github.com/lsst/qserv-operator/pkg/constants"
)

// GetCzarName returns the name for czar ressources
func GetCzarName(r *qservv1alpha1.Qserv) string {
	return generateName(r.Name, constants.CzarName)
}

// GetWorkerName returns the name for xrootd redirector ressources
func GetWorkerName(r *qservv1alpha1.Qserv) string {
	return generateName(r.Name, constants.WorkerName)
}

// GetXrootdRedirectorName returns the name for xrootd redirector ressources
func GetXrootdRedirectorName(r *qservv1alpha1.Qserv) string {
	return generateName(r.Name, constants.XrootdRedirectorName)
}

// // GetRedisShutdownName returns the name for redis resources
// func GetRedisShutdownName(r *qservv1alpha1.Qserv) string {
//         return generateName(constants.RedisShutdownName, r.Name)
// }

// GetXrootdName returns the name for xrootd resources
func GetXrootdName(r *qservv1alpha1.Qserv) string {
	return generateName(r.Name, constants.XrootdName)
}

func generateName(metaName, typeName string) string {
	return fmt.Sprintf("%s-%s", metaName, typeName)
}
