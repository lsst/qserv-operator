package specs

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/go-logr/logr"
	qservv1beta1 "github.com/lsst/qserv-operator/api/v1beta1"
	"github.com/lsst/qserv-operator/controllers/constants"
	"github.com/lsst/qserv-operator/controllers/util"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type templateData struct {
	CzarDatabaseRootURL                string
	MariadbSocket                      string
	QservInstance                      string
	QstatusMysqldHost                  string
	HTTPPort                           uint
	ReplicationControllerFQDN          string // Example: qserv-repl-ctl-0.qserv-repl-ctl.default.svc.cluster.local
	ReplicationDatabaseURL             string
	ReplicationDatabaseRootURL         string
	ReplicationRegistryDN              string
	ReplicationLoaderProcessingThreads uint
	ResultsProtocol                    qservv1beta1.ResultsProtocolType
	SocketQservUser                    string
	SocketRootUser                     string
	WmgrPort                           uint
	WorkerDatabaseLocalRootURL         string
	WorkerDatabaseLocalURL             string
	WorkerIngestNumRetries             uint
	WorkerIngestMaxRetries             uint
	WorkerDN                           string
	WorkerReplicas                     uint
	XrootdRedirectorDN                 string
	XrootdRedirectorReplicas           uint
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

	// TODO see https://medium.com/@etiennerouzeaud/change-default-delimiters-templating-with-template-delims-857938a0b661
	// to change template delimiters
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
	fileinfo, err := os.Stat(root)
	if os.IsNotExist(err) {
		reqLogger.V(0).Info("Path does not exist", "path", root)
	} else if !fileinfo.IsDir() {
		err = fmt.Errorf("path is not a directory: %s", root)
	} else {
		err = filepath.Walk(root,
			func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					reqLogger.V(5).Info("Scan path", "path", path)
					files[info.Name()], _ = applyTemplate(path, tmplData)
				}
				return nil
			})
	}
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
		loaderProcessingThreads = constants.ReplicationWorkerThreadFactor * limit
	}
	return loaderProcessingThreads
}

func generateTemplateData(r *qservv1beta1.Qserv) templateData {
	cpuLimit := r.Spec.Worker.ReplicationResources.Limits.Cpu()
	return templateData{
		CzarDatabaseRootURL:                util.GetCzarDatabaseRootURL(r),
		MariadbSocket:                      constants.MariadbSocket,
		QservInstance:                      util.GetQservInstanceName(r),
		QstatusMysqldHost:                  util.GetCzarServiceName(r),
		HTTPPort:                           constants.HTTPPort,
		ReplicationControllerFQDN:          util.GetReplCtlFQDN(r),
		ReplicationDatabaseURL:             util.GetReplicationDatabaseURL(r),
		ReplicationDatabaseRootURL:         util.GetReplicationDatabaseRootURL(r),
		ReplicationRegistryDN:              util.GetReplicationRegistryServiceName(r),
		ResultsProtocol:                    r.Spec.Worker.ResultsProtocol,
		SocketQservUser:                    util.SocketQservUser,
		SocketRootUser:                     util.SocketRootUser,
		WorkerDN:                           util.GetWorkerServiceName(r),
		WorkerDatabaseLocalRootURL:         util.WorkerDatabaseLocalRootURL,
		WorkerDatabaseLocalURL:             util.WorkerDatabaseLocalURL,
		WorkerIngestNumRetries:             constants.WorkerIngestNumRetries,
		WorkerIngestMaxRetries:             constants.WorkerIngestMaxRetries,
		WorkerReplicas:                     uint(r.Spec.Worker.Replicas),
		XrootdRedirectorDN:                 util.GetXrootdRedirectorServiceName(r),
		XrootdRedirectorReplicas:           uint(r.Spec.Xrootd.Replicas),
		ReplicationLoaderProcessingThreads: getReplicationWorkerThread(cpuLimit),
	}
}

// ConfigMapSpec provide default procedures for all Qserv configmaps specifications
type ConfigMapSpec struct {
	qserv *qservv1beta1.Qserv
}

// Initialize initialize configmap specification
func (c *ConfigMapSpec) Initialize(qserv *qservv1beta1.Qserv) client.Object {
	c.qserv = qserv
	var object client.Object = &v1.ConfigMap{}
	return object
}

// Update update configmap specification for Qserv containers
func (c *ConfigMapSpec) Update(object client.Object) (bool, error) {
	return false, nil
}

// ContainerConfigMapSpec provide procedures for all Qserv containers configmaps specifications
type ContainerConfigMapSpec struct {
	ConfigMapSpec
	ContainerName constants.ContainerName
	Subdir        string
}

// GetName return name for container ConfigMaps
func (c *ContainerConfigMapSpec) GetName() string {
	suffix := fmt.Sprintf("%s-%s", c.ContainerName, c.Subdir)
	return util.PrefixConfigmap(c.qserv, suffix)
}

// Create can generate 2 kind of configmaps for Qserv containers
// one with startup scripts and one with configuration files
func (c *ContainerConfigMapSpec) Create() (client.Object, error) {
	cr := c.qserv
	tmplData := generateTemplateData(cr)

	reqLogger := log.WithValues("Request.Namespace", cr.Namespace, "Request.Name", cr.Name)

	name := c.GetName()
	namespace := cr.Namespace

	labels := util.GetContainerLabels(c.ContainerName, cr.Name)
	root := filepath.Join("/", "configmap", string(c.ContainerName), c.Subdir)

	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Data: scanDir(root, reqLogger, &tmplData),
	}
	return cm, nil
}

// SQLConfigMapSpec provide procedures for all Qserv ConfigMaps related to database initialization
type SQLConfigMapSpec struct {
	ConfigMapSpec
	Database constants.PodClass
}

// GetName return name for SQL ConfigMaps
func (c *SQLConfigMapSpec) GetName() string {
	suffix := fmt.Sprintf("sql-%s", c.Database)
	return util.PrefixConfigmap(c.qserv, suffix)
}

// Create generate ConfigMaps for initContainers in charge of databases initializations
func (c *SQLConfigMapSpec) Create() (client.Object, error) {
	cr := c.qserv
	tmplData := generateTemplateData(cr)

	reqLogger := log.WithValues("Request.Namespace", cr.Namespace, "Request.Name", cr.Name)
	// reqLogger.Info("XXXXX %s", "tmplData", tmplData)

	name := c.GetName()
	namespace := cr.Namespace

	labels := util.GetComponentLabels(c.Database, cr.Name)
	root := filepath.Join("/", "configmap", "initdb", "sql", string(c.Database))

	configmap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Data: scanDir(root, reqLogger, &tmplData),
	}
	return configmap, nil
}

// DotQservConfigMapSpec provide procedures for .qserv ConfigMap
type DotQservConfigMapSpec struct {
	ConfigMapSpec
}

// GetName return name for .qserv ConfigMap
func (c *DotQservConfigMapSpec) GetName() string {
	return util.PrefixConfigmap(c.qserv, constants.DotQserv)
}

// Create generate configmap for Qserv client configuration
func (c *DotQservConfigMapSpec) Create() (client.Object, error) {
	cr := c.qserv
	tmplData := generateTemplateData(cr)

	reqLogger := log.WithValues("Request.Namespace", cr.Namespace, "Request.Name", cr.Name)

	name := c.GetName()
	namespace := cr.Namespace

	labels := util.GetComponentLabels(constants.Czar, cr.Name)
	root := filepath.Join("/", "configmap", "dot-qserv")

	cm := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Data: scanDir(root, reqLogger, &tmplData),
	}

	return cm, nil
}
