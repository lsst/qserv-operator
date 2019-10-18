package sync

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	qservv1alpha1 "github.com/lsst/qserv-operator/pkg/apis/qserv/v1alpha1"
	"github.com/lsst/qserv-operator/pkg/constants"
	"github.com/lsst/qserv-operator/pkg/scheme/qserv"
	"github.com/lsst/qserv-operator/pkg/staging/syncer"
)

// NewSecretSyncer returns a new sync.Interface for reconciling Qserv secrets
func NewSecretSyncer(r *qservv1alpha1.Qserv, c client.Client, scheme *runtime.Scheme, service constants.ContainerName) syncer.Interface {
	cm := qserv.GenerateSecret(r, controllerLabels, service)
	objectName := fmt.Sprintf("%sSecret", strings.Title(string(service)))
	return syncer.NewObjectSyncer(objectName, r, cm, c, scheme, func(existing runtime.Object) error {
		return nil
	})
}
