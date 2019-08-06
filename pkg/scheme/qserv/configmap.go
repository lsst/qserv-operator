package qserv

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-logr/logr"
	qservv1alpha1 "github.com/lsst/qserv-operator/pkg/apis/qserv/v1alpha1"
	"github.com/lsst/qserv-operator/pkg/constants"
	"github.com/lsst/qserv-operator/pkg/util"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type filedesc struct {
	name    string
	content []byte
}

func getFileContent(path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Error(err, fmt.Sprintf("Cannot open file: %s", path))
		os.Exit(1)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error(err, fmt.Sprintf("Cannot read file: %s", path))
		os.Exit(1)
	}
	return fmt.Sprintf("%s", b)
}

// TODO manage secret cleanly
func getSecretData(r *qservv1alpha1.Qserv, service string) map[string][]byte {
	files := make(map[string][]byte)
	if service == "mariadb" {
		files["mariadb.secret.sh"] = []byte(`MYSQL_ROOT_PASSWORD="CHANGEME"
		MYSQL_MONITOR_PASSWORD="CHANGEMETOO"`)
	} else if service == "wmgr" {
		files["wmgr.secret"] = []byte(`USER:CHANGEMEAGAIN`)
	}
	return files
}

func GenerateSecret(r *qservv1alpha1.Qserv, labels map[string]string, service string) *v1.Secret {
	name := fmt.Sprintf("secret-%s", service)
	namespace := r.Namespace

	labels = util.MergeLabels(labels, util.GetLabels(constants.XrootdName, r.Name))

	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Data: getSecretData(r, service),
	}
}

func scanDir(root string, reqLogger logr.Logger) map[string]string {
	files := make(map[string]string)
	reqLogger.Info(fmt.Sprintf("Walk through %s", root))
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			reqLogger.Info(fmt.Sprintf("Scan %s", path))
			files[info.Name()] = getFileContent(path)
		}
		return nil
	})
	if err != nil {
		reqLogger.Error(err, fmt.Sprintf("Cannot walk path: %s", root))
		os.Exit(1)
	}
	return files
}

func GenerateMicroserviceConfigMap(r *qservv1alpha1.Qserv, labels map[string]string, service string, subdir string) *v1.ConfigMap {
	reqLogger := log.WithValues("Request.Namespace", r.Namespace, "Request.Name", r.Name)

	name := fmt.Sprintf("config-%s-%s", service, subdir)
	namespace := r.Namespace

	labels = util.MergeLabels(labels, util.GetLabels(constants.XrootdName, r.Name))
	root := filepath.Join("/", "configmap", service, subdir)

	return &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Data: scanDir(root, reqLogger),
	}
}

func GenerateSqlConfigMap(r *qservv1alpha1.Qserv, labels map[string]string, db string) *v1.ConfigMap {
	reqLogger := log.WithValues("Request.Namespace", r.Namespace, "Request.Name", r.Name)

	name := fmt.Sprintf("config-sql-%s", db)
	namespace := r.Namespace

	labels = util.MergeLabels(labels, util.GetLabels(constants.XrootdName, r.Name))
	root := filepath.Join("/", "configmap", "init", "sql", db)

	return &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Data: scanDir(root, reqLogger),
	}
}

func GenerateDotQservConfigMap(r *qservv1alpha1.Qserv, labels map[string]string) *v1.ConfigMap {
	reqLogger := log.WithValues("Request.Namespace", r.Namespace, "Request.Name", r.Name)

	name := "config-dot-qserv"
	namespace := r.Namespace

	labels = util.MergeLabels(labels, util.GetLabels(constants.XrootdName, r.Name))
	root := filepath.Join("/", "configmap", "dot-qserv")

	return &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Data: scanDir(root, reqLogger),
	}
}
