/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/go-logr/logr"
	qservv1beta1 "github.com/lsst/qserv-operator/api/v1beta1"
	"github.com/lsst/qserv-operator/pkg/constants"
	"github.com/lsst/qserv-operator/pkg/syncer"
	"github.com/lsst/qserv-operator/pkg/syncers"
	appsv1 "k8s.io/api/apps/v1"
)

// QservReconciler reconciles a Qserv object
type QservReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=qserv.lsst.org,resources=qservs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=qserv.lsst.org,resources=qservs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=qserv.lsst.org,resources=qservs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Qserv object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.9.2/pkg/reconcile
func (r *QservReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	// TODO check which log to use
	// log := r.Log.WithValues("qserv", request.NamespacedName)

	log.Info("Reconciling Qserv")

	// Fetch the Qserv instance
	qserv := &qservv1beta1.Qserv{}
	err := r.Get(ctx, request.NamespacedName, qserv)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("Qserv resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get Qserv")
		return ctrl.Result{}, err
	}

	// Manage default values for specification
	r.Scheme.Default(qserv)

	result, err := r.updateQservStatus(ctx, request, qserv, &log)
	if err != nil {
		log.Error(err, "Unable to update Qserv status")
		return result, err
	}

	// Manage syncronisation
	qservSyncers := []syncer.Interface{
		syncers.NewCzarStatefulSetSyncer(qserv, r.Client, r.Scheme),
		syncers.NewDotQservConfigMapSyncer(qserv, r.Client, r.Scheme),
		syncers.NewDashboardDeploymentSyncer(qserv, r.Client, r.Scheme),
		syncers.NewDashboardServiceSyncer(qserv, r.Client, r.Scheme),
		syncers.NewWorkerStatefulSetSyncer(qserv, r.Client, r.Scheme),
		syncers.NewReplicationCtlServiceSyncer(qserv, r.Client, r.Scheme),
		syncers.NewReplicationCtlStatefulSetSyncer(qserv, r.Client, r.Scheme),
		syncers.NewIngestDbServiceSyncer(qserv, r.Client, r.Scheme),
		syncers.NewIngestDbStatefulSetSyncer(qserv, r.Client, r.Scheme),
		syncers.NewReplicationDbServiceSyncer(qserv, r.Client, r.Scheme),
		syncers.NewReplicationDbStatefulSetSyncer(qserv, r.Client, r.Scheme),
		syncers.NewXrootdRedirectorServiceSyncer(qserv, r.Client, r.Scheme),
		syncers.NewXrootdStatefulSetSyncer(qserv, r.Client, r.Scheme),
	}

	qservSyncers = append(syncers.NewQservServicesSyncer(qserv, r.Client, r.Scheme), qservSyncers...)

	for _, configmapClass := range constants.ContainerConfigmaps {
		for _, subpath := range []string{"etc", "start"} {
			qservSyncers = append(qservSyncers, syncers.NewContainerConfigMapSyncer(qserv, r.Client, r.Scheme, configmapClass, subpath))
		}
	}
	qservSyncers = append(qservSyncers, syncers.NewContainerConfigMapSyncer(qserv, r.Client, r.Scheme, constants.InitDbName, "start"))

	for _, db := range constants.Databases {
		qservSyncers = append(qservSyncers, syncers.NewSQLConfigMapSyncer(qserv, r.Client, r.Scheme, db))
	}

	// Specify Network Policies
	if qserv.Spec.NetworkPolicies {
		qservSyncers = append(qservSyncers, syncers.NewNetworkPoliciesSyncer(qserv, r.Client, r.Scheme)...)
	}

	if err = r.sync(qservSyncers); err != nil {
		return reconcile.Result{}, err
	}

	// Update status.Nodes if needed
	/* 	if !reflect.DeepEqual(podNames, qserv.Status.Nodes) {
		qserv.Status.Nodes = podNames
		err := r.Status().Update(ctx, qserv)
		if err != nil {
			log.Error(err, "Failed to update Qserv status")
			return ctrl.Result{}, err
		}
	} */

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *QservReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&qservv1beta1.Qserv{}).
		Owns(&appsv1.StatefulSet{}).
		Complete(r)
}
