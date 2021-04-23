package qserv

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/go-logr/logr"
	qservv1alpha1 "github.com/lsst/qserv-operator/api/v1alpha1"
	"github.com/lsst/qserv-operator/pkg/constants"
	"github.com/lsst/qserv-operator/pkg/util"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type templateData struct {
	CzarDomainName            string
	DashboardDn               string
	DashboardPort             uint
	QstatusMysqldHost         string
	ReplicationControllerPort uint
	// Example: qserv-repl-ctl-0.qserv-repl-ctl.default.svc.cluster.local
	ReplicationControllerFQDN          string
	ReplicationLoaderProcessingThreads uint
	WmgrPort                           uint
	WorkerDn                           string
	WorkerReplicas                     uint
	XrootdRedirectorDn                 string
	XrootdReplicas                     uint
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func applyTemplate(path string, tmplData *templateData) (string, error) {

	if !fileExists(path) {
		return "", fmt.Errorf("file does not exists: %s", path)
	}

	tmpl, err := template.New(filepath.Base(path)).Funcs(util.TemplateFunctions).ParseFiles(path)
	if err != nil {
		log.Error(err, fmt.Sprintf("cannot open template file: %s", path))
		return "", nil
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, tmplData)

	if err != nil {
		log.Error(err, fmt.Sprintf("Cannot apply template: %s", path))
	}
	return buf.String(), nil
}

func scanDir(root string, reqLogger logr.Logger, tmplData *templateData) map[string]string {
	files := make(map[string]string)
	reqLogger.Info(fmt.Sprintf("Walk through %s", root))
	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				reqLogger.Info(fmt.Sprintf("Scan %s", path))
				files[info.Name()], _ = applyTemplate(path, tmplData)
			}
			return nil
		})
	if err != nil {
		reqLogger.Error(err, fmt.Sprintf("Cannot walk path: %s", root))
	}
	return files
}

func getReplicationWorkerThread(cpuLimit *resource.Quantity) uint {
	var loaderProcessingThreads uint
	limit := uint(cpuLimit.Value())
	if limit == 0 {
		loaderProcessingThreads = constants.ReplicationWorkerDefaultThreads
	} else {
		loaderProcessingThreads = 2 * limit
	}
	return loaderProcessingThreads
}

func generateTemplateData(r *qservv1alpha1.Qserv) templateData {
	cpuLimit := r.Spec.Worker.ReplicationResources.Limits.Cpu()
	return templateData{
		CzarDomainName:                     util.GetCzarServiceName(r),
		DashboardDn:                        util.GetDashboardServiceName(r),
		DashboardPort:                      constants.DashboardPort,
		QstatusMysqldHost:                  util.GetCzarServiceName(r),
		ReplicationControllerPort:          constants.ReplicationControllerPort,
		ReplicationControllerFQDN:          util.GetReplCtlFQDN(r),
		WorkerDn:                           util.GetWorkerServiceName(r),
		WmgrPort:                           constants.WmgrPort,
		WorkerReplicas:                     uint(r.Spec.Worker.Replicas),
		XrootdRedirectorDn:                 util.GetXrootdRedirectorServiceName(r),
		XrootdReplicas:                     uint(r.Spec.Xrootd.Replicas),
		ReplicationLoaderProcessingThreads: getReplicationWorkerThread(cpuLimit),
	}
}

// GenerateContainerConfigMap generate 2 configmaps for Qserv containers
// one with startup scripts and one with configuration files
func GenerateContainerConfigMap(r *qservv1alpha1.Qserv, container constants.ContainerName, subdir string) *v1.ConfigMap {

	tmplData := generateTemplateData(r)

	reqLogger := log.WithValues("Request.Namespace", r.Namespace, "Request.Name", r.Name)

	name := util.PrefixConfigmap(r, fmt.Sprintf("%s-%s", container, subdir))
	namespace := r.Namespace

	labels := util.GetContainerLabels(container, r.Name)
	root := filepath.Join("/", "configmap", string(container), subdir)

	return &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Data: scanDir(root, reqLogger, &tmplData),
	}
}

// GenerateSQLConfigMap generate configmaps for initContainers in charge of databases initializations
func GenerateSQLConfigMap(r *qservv1alpha1.Qserv, db constants.PodClass) *v1.ConfigMap {

	tmplData := generateTemplateData(r)

	reqLogger := log.WithValues("Request.Namespace", r.Namespace, "Request.Name", r.Name)
	// reqLogger.Info("XXXXX %s", "tmplData", tmplData)

	name := util.PrefixConfigmap(r, fmt.Sprintf("sql-%s", db))
	namespace := r.Namespace

	labels := util.GetComponentLabels(db, r.Name)
	root := filepath.Join("/", "configmap", "initdb", "sql", string(db))

	return &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Data: scanDir(root, reqLogger, &tmplData),
	}
}

// GenerateDotQservConfigMap generate configmap for Qserv client configuration
func GenerateDotQservConfigMap(cr *qservv1alpha1.Qserv) *v1.ConfigMap {

	tmplData := generateTemplateData(cr)

	reqLogger := log.WithValues("Request.Namespace", cr.Namespace, "Request.Name", cr.Name)

	name := util.PrefixConfigmap(cr, "dot-qserv")
	namespace := cr.Namespace

	labels := util.GetComponentLabels(constants.Czar, cr.Name)
	root := filepath.Join("/", "configmap", "dot-qserv")

	return &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Data: scanDir(root, reqLogger, &tmplData),
	}
}
