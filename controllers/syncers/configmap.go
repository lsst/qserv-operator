package syncers

import (
	"fmt"
	"strings"

	"github.com/lsst/qserv-operator/controllers/constants"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	qservv1beta1 "github.com/lsst/qserv-operator/api/v1beta1"
	"github.com/lsst/qserv-operator/controllers/objects"
	"github.com/lsst/qserv-operator/controllers/syncer"
)

// NewContainerConfigMapSyncer generate configmap specification for a given container
func NewContainerConfigMapSyncer(r *qservv1beta1.Qserv, c client.Client, scheme *runtime.Scheme, container constants.ContainerName, subpath string) syncer.Interface {
	cm := objects.GenerateContainerConfigMap(r, container, subpath)
	objectName := fmt.Sprintf("%s%sConfigMap", strings.Title(string(container)), strings.Title(subpath))
	return syncer.NewObjectSyncer(objectName, r, cm, c, scheme, func() error {
		return nil
	})
}

// NewDotQservConfigMapSyncer generate configmap specification for Qserv clients
func NewDotQservConfigMapSyncer(r *qservv1beta1.Qserv, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	cm := objects.GenerateDotQservConfigMap(r)
	return syncer.NewObjectSyncer("DotQservConfigMap", r, cm, c, scheme, func() error {
		return nil
	})
}

// NewSQLConfigMapSyncer generate configmap specification for initContainer in charge of database initialization
func NewSQLConfigMapSyncer(r *qservv1beta1.Qserv, c client.Client, scheme *runtime.Scheme, db constants.PodClass) syncer.Interface {
	cm := objects.GenerateSQLConfigMap(r, db)
	objectName := fmt.Sprintf("%sSqlConfigMap", strings.Title(string(db)))
	return syncer.NewObjectSyncer(objectName, r, cm, c, scheme, func() error {
		return nil
	})
}
