package objects

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
	CzarDomainName            string
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

func generateTemplateData(r *qservv1beta1.Qserv) templateData {
	cpuLimit := r.Spec.Worker.ReplicationResources.Limits.Cpu()
	return templateData{
		CzarDomainName:                     util.GetCzarServiceName(r),
		QstatusMysqldHost:                  util.GetCzarServiceName(r),
		ReplicationControllerPort:          constants.ReplicationControllerPort,
		ReplicationControllerFQDN:          util.GetReplCtlFQDN(r),
		WorkerDn:                           util.GetWorkerServiceName(r),
		WorkerReplicas:                     uint(r.Spec.Worker.Replicas),
		XrootdRedirectorDn:                 util.GetXrootdRedirectorServiceName(r),
		XrootdReplicas:                     uint(r.Spec.Xrootd.Replicas),
		ReplicationLoaderProcessingThreads: getReplicationWorkerThread(cpuLimit),
	}
}

type ContainerConfigMapSpec struct {
	qserv         *qservv1beta1.Qserv
	ContainerName constants.ContainerName
	Subdir        string
}

func (c *ContainerConfigMapSpec) Initialize(qserv *qservv1beta1.Qserv) client.Object {
	c.qserv = qserv
	var object client.Object = &v1.ConfigMap{}
	return object
}

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

// Update update configmap specification for Qserv containers
func (c *ContainerConfigMapSpec) Update(object client.Object) (bool, error) {
	return false, nil
}

type SQLConfigMapSpec struct {
	qserv    *qservv1beta1.Qserv
	Database constants.PodClass
}

func (c *SQLConfigMapSpec) Initialize(qserv *qservv1beta1.Qserv) client.Object {
	c.qserv = qserv
	var object client.Object = &v1.ConfigMap{}
	return object
}

func (c *SQLConfigMapSpec) GetName() string {
	suffix := fmt.Sprintf("%s-%s", c.ContainerName, c.Subdir)
	return util.PrefixConfigmap(c.qserv, suffix)
}

// Create generate configmaps for initContainers in charge of databases initializations
func (c *SQLConfigMapSpec) Create() (client.Object, error) {
	cr := c.qserv
	tmplData := generateTemplateData(cr)

	reqLogger := log.WithValues("Request.Namespace", cr.Namespace, "Request.Name", cr.Name)
	// reqLogger.Info("XXXXX %s", "tmplData", tmplData)

	name := util.PrefixConfigmap(r, fmt.Sprintf("sql-%s", c.Database))
	namespace := r.Namespace

	labels := util.GetComponentLabels(db, r.Name)
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

// Update update configmap for Qserv client configuration
func (c *SQLConfigMapSpec) Update(object client.Object) (bool, error) {
	return false, nil
}

type DotQservConfigMapSpec struct {
	qserv *qservv1beta1.Qserv
}

func (c *DotQservConfigMapSpec) GetName() string {
	return util.PrefixConfigmap(c.qserv, constants.DotQserv)
}

func (c *DotQservConfigMapSpec) Initialize(qserv *qservv1beta1.Qserv) client.Object {
	c.qserv = qserv
	var object client.Object = &v1.ConfigMap{}
	return object
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

// Update update configmap for Qserv client configuration
func (c *DotQservConfigMapSpec) Update(object client.Object) (bool, error) {
	return false, nil
}
