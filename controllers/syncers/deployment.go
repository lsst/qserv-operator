package syncers

import (
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	qservv1beta1 "github.com/lsst/qserv-operator/api/v1beta1"
	"github.com/lsst/qserv-operator/controllers/objects"
	"github.com/lsst/qserv-operator/controllers/syncer"
	"github.com/lsst/qserv-operator/controllers/util"
)

// NewDashboardDeploymentSyncer returns a new sync.Interface for reconciling Qserv Dashboard Deployment
func NewDashboardDeploymentSyncer(r *qservv1beta1.Qserv, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	deployment := objects.GenerateDashboardDeployment(r)
	return syncer.NewObjectSyncer("DashboardDeployment", r, deployment, c, scheme, util.NoFunc)
}
