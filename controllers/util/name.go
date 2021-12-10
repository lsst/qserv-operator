package util

import (
	"fmt"
	"net"
	"strings"

	qservv1beta1 "github.com/lsst/qserv-operator/api/v1beta1"
	"github.com/lsst/qserv-operator/controllers/constants"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("name")

// GetName returns a name whose prefix is instance name and suffix typeName
func GetName(r *qservv1beta1.Qserv, typeName string) string {
	return fmt.Sprintf("%s-%s", r.Name, typeName)
}

// GetCzarServiceName returns name of Qserv czar headless service
func GetCzarServiceName(cr *qservv1beta1.Qserv) string {
	return GetName(cr, string(constants.Czar))
}

// GetWorkerServiceName returns name of Qserv workers headless service
func GetWorkerServiceName(cr *qservv1beta1.Qserv) string {
	return GetName(cr, string(constants.Worker))
}

// GetReplCtlServiceName returns name of Replication Con headless service
func GetReplCtlServiceName(cr *qservv1beta1.Qserv) string {
	return GetName(cr, string(constants.ReplCtl))
}

// GetWorkerNameFilter returns a filter on hostname for mysql user
// Example: use in "CREATE USER 'qsreplica'@'<filter>'"
func GetWorkerNameFilter(cr *qservv1beta1.Qserv) string {
	filter := cr.Name + "-" + string(constants.Worker) + "-%." + GetWorkerServiceName(cr) + "." + cr.GetNamespace() + ".svc." + getClusterDomain()
	return filter
}

// GetReplCtlFQDN returns a Replication Controller FQDN
// It can be used for mysql authentication
// Example: use in "CREATE USER 'qsreplica'@'<FQDN>'"
func GetReplCtlFQDN(cr *qservv1beta1.Qserv) string {
	fqdn := cr.Name + "-" + string(constants.ReplCtlName) + "-0." + GetReplCtlServiceName(cr) + "." + cr.GetNamespace() + ".svc." + getClusterDomain()
	return fqdn
}

// GetXrootdRedirectorServiceName returns name of Xrootd redirector headless service
func GetXrootdRedirectorServiceName(cr *qservv1beta1.Qserv) string {
	return GetName(cr, string(constants.XrootdRedirector))
}

// GetClusterDomain returns Kubernetes cluster domain, default to "cluster.local"
func getClusterDomain() string {
	apiSvc := "kubernetes.default.svc"

	defaultClusterDomain := "cluster.local"

	cname, err := net.LookupCNAME(apiSvc)
	if err != nil {
		log.V(2).Info("Unable to get fqdn for %v, using '%v'", defaultClusterDomain)
		return defaultClusterDomain
	}

	clusterDomain := strings.TrimPrefix(cname, apiSvc)
	clusterDomain = strings.TrimPrefix(clusterDomain, ".")
	clusterDomain = strings.TrimSuffix(clusterDomain, ".")

	return clusterDomain
}

// PrefixConfigmap add a common prefix to all ConfigMap names of a given Qserv instance
func PrefixConfigmap(r *qservv1beta1.Qserv, name string) string {
	return fmt.Sprintf("%s-%s", r.Name, name)
}

// GetConfigVolumeName add a common prefix to all pod volume names attaching a configmap
func GetConfigVolumeName(suffix string) string {
	return fmt.Sprintf("config-%s", suffix)
}

// GetSecretName return the name of a secret for a given container and a given Qserv instance
func GetSecretName(cr *qservv1beta1.Qserv, containerName constants.ContainerName) string {
	return fmt.Sprintf("%s-%s", GetSecretVolumeName(containerName), cr.Name)
}

// GetSecretVolumeName return the name of a volume for a given secret
func GetSecretVolumeName(containerName constants.ContainerName) string {
	return fmt.Sprintf("secret-%s", containerName)
}
