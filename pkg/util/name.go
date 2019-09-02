package util

import (
	"fmt"

	qservv1alpha1 "github.com/lsst/qserv-operator/pkg/apis/qserv/v1alpha1"
	"github.com/lsst/qserv-operator/pkg/constants"
)

// GetCzarName returns the name for Czar ressources
func GetCzarName(r *qservv1alpha1.Qserv) string {
	return generateName(r.Name, constants.CzarName)
}

// GetReplicationCtlName returns the name for Replication Controller ressources
func GetReplicationCtlName(r *qservv1alpha1.Qserv) string {
	return generateName(r.Name, constants.ReplDbName)
}

// GetReplicationDbName returns the name for Replication Db ressources
func GetReplicationDbName(r *qservv1alpha1.Qserv) string {
	return generateName(r.Name, constants.ReplDbName)
}

// GetWorkerName returns the name for Xrootd redirector ressources
func GetWorkerName(r *qservv1alpha1.Qserv) string {
	return generateName(r.Name, constants.WorkerName)
}

// GetXrootdRedirectorName returns the name for Xrootd redirector ressources
func GetXrootdRedirectorName(r *qservv1alpha1.Qserv) string {
	return generateName(r.Name, constants.XrootdRedirectorName)
}

// GetXrootdName returns the name for Xrootd resources
func GetXrootdName(r *qservv1alpha1.Qserv) string {
	return generateName(r.Name, constants.XrootdName)
}

func generateName(metaName, typeName string) string {
	return fmt.Sprintf("%s-%s", metaName, typeName)
}
