package sync

import (
	"fmt"
	"strings"

	"github.com/lsst/qserv-operator/pkg/constants"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	qservv1alpha1 "github.com/lsst/qserv-operator/pkg/apis/qserv/v1alpha1"
	"github.com/lsst/qserv-operator/pkg/scheme/qserv"
	"github.com/lsst/qserv-operator/pkg/staging/syncer"
)

// NewContainerConfigMapSyncer generate configmap specification for a given container
func NewContainerConfigMapSyncer(r *qservv1alpha1.Qserv, c client.Client, scheme *runtime.Scheme, container constants.ContainerName, subpath string) syncer.Interface {
	cm := qserv.GenerateContainerConfigMap(r, controllerLabels, container, subpath)
	objectName := fmt.Sprintf("%s%sConfigMap", strings.Title(string(container)), strings.Title(subpath))
	return syncer.NewObjectSyncer(objectName, r, cm, c, scheme, func() error {
		return nil
	})
}

func NewDotQservConfigMapSyncer(r *qservv1alpha1.Qserv, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	cm := qserv.GenerateDotQservConfigMap(r, controllerLabels)
	return syncer.NewObjectSyncer("DotQservConfigMap", r, cm, c, scheme, func() error {
		return nil
	})
}

func NewSqlConfigMapSyncer(r *qservv1alpha1.Qserv, c client.Client, scheme *runtime.Scheme, db constants.ComponentName) syncer.Interface {
	cm := qserv.GenerateSqlConfigMap(r, controllerLabels, db)
	objectName := fmt.Sprintf("%sSqlConfigMap", strings.Title(string(db)))
	return syncer.NewObjectSyncer(objectName, r, cm, c, scheme, func() error {
		return nil
	})
}
