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

// NewSQLConfigMapSyncer generate configmap specification for initContainer in charge of database initialization
func NewSQLConfigMapSyncer(r *qservv1beta1.Qserv, c client.Client, scheme *runtime.Scheme, db constants.PodClass) syncer.Interface {
	cm := objects.GenerateSQLConfigMap(r, db)
	objectName := fmt.Sprintf("%sSqlConfigMap", strings.Title(string(db)))
	return syncer.NewObjectSyncer(objectName, r, cm, c, scheme, func() error {
		return nil
	})
}
