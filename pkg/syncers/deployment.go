package syncers

import (
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	qservv1alpha1 "github.com/lsst/qserv-operator/api/v1alpha1"
	qserv "github.com/lsst/qserv-operator/pkg/resources"
	"github.com/lsst/qserv-operator/pkg/syncer"
)

// NewDashboardDeploymentSyncer returns a new sync.Interface for reconciling Qserv Dashboard Deployment
func NewDashboardDeploymentSyncer(r *qservv1alpha1.Qserv, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	deployment := qserv.GenerateDashboardDeployment(r, controllerLabels)
	return syncer.NewObjectSyncer("DashboardDeployment", r, deployment, c, scheme, noFunc)
}
