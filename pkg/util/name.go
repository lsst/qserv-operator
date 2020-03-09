package util

import (
	"fmt"
	"net"
	"strings"

	qservv1alpha1 "github.com/lsst/qserv-operator/pkg/apis/qserv/v1alpha1"
	"github.com/lsst/qserv-operator/pkg/constants"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("name")

// GetName returns a name whose prefix is instance name and suffix typeName
func GetName(r *qservv1alpha1.Qserv, typeName string) string {
	return fmt.Sprintf("%s-%s", r.Name, typeName)
}

// GetCzarServiceName returns name of Qserv czar headless service
func GetCzarServiceName(cr *qservv1alpha1.Qserv) string {
	return GetName(cr, string(constants.CzarName))
}

// GetWorkerServiceName returns name of Qserv workers headless service
func GetWorkerServiceName(cr *qservv1alpha1.Qserv) string {
	return GetName(cr, string(constants.WorkerName))
}

// GetReplCtlServiceName returns name of Replication Con headless service
func GetReplCtlServiceName(cr *qservv1alpha1.Qserv) string {
	return GetName(cr, string(constants.ReplCtlName))
}

// GetWorkerNameFilter returns a filter on hostname for mysql user
// Example: use in "CREATE USER 'qsreplica'@'<filter>'"
func GetWorkerNameFilter(cr *qservv1alpha1.Qserv) string {
	filter := cr.Name + "-" + string(constants.WorkerName) + "-%." + GetWorkerServiceName(cr) + "." + cr.GetNamespace() + ".svc" + getClusterDomain()
	return filter
}

// GetReplCtlNameFilter returns a filter on hostname for mysql user
// Example: use in "CREATE USER 'qsreplica'@'<filter>'"
func GetReplCtlNameFilter(cr *qservv1alpha1.Qserv) string {
	filter := cr.Name + "-" + string(constants.ReplCtlName) + "-%." + GetReplCtlServiceName(cr) + "." + cr.GetNamespace() + ".svc" + getClusterDomain()
	return filter
}

// GetXrootdRedirectorServiceName returns name of Xrootd redirector headless service
func GetXrootdRedirectorServiceName(cr *qservv1alpha1.Qserv) string {
	return GetName(cr, string(constants.XrootdRedirectorName))
}

// GetClusterDomain returns Kubernetes cluster domain, default to "cluster.local"
func getClusterDomain() string {
	apiSvc := "kubernetes.default.svc"

	clusterDomain := "cluster.local"

	cname, err := net.LookupCNAME(apiSvc)
	if err != nil {
		log.V(2).Info("Unable to get fqdn for %v, using '%v'", clusterDomain)
		return clusterDomain
	}

	clusterDomain = strings.TrimPrefix(cname, apiSvc)
	clusterDomain = strings.TrimSuffix(clusterDomain, ".")

	return clusterDomain
}

// PrefixConfigmap add a common prefix to all ConfigMap names of a given Qserv instance
func PrefixConfigmap(r *qservv1alpha1.Qserv, name string) string {
	return fmt.Sprintf("%s-%s", r.Name, name)
}

// PrefixConfigMap add a common prefix to all ConfigMap names of a given Qserv instance
func GetConfigVolumeName(suffix string) string {
	return fmt.Sprintf("config-%s", suffix)
}

// GetSecretName return the name of a secret for a given container and a given Qserv instance
func GetSecretName(cr *qservv1alpha1.Qserv, containerName constants.ContainerName) string {
	return fmt.Sprintf("%s-%s", GetSecretVolumeName(containerName), cr.Name)
}

// GetSecretVolumeName return the name of a volume for a given secret
func GetSecretVolumeName(containerName constants.ContainerName) string {
	return fmt.Sprintf("secret-%s", containerName)
}

// GetXrootdReplicas returns num of Xrootd replicas defined in cr.Spec.Xrootd.Replicas
func GetXrootdReplicas(cr *qservv1alpha1.Qserv) string {
	return (cr, cr.Spec.Xrootd.Replicas)
}
