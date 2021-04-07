/*

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
	"fmt"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	qservv1alpha1 "github.com/lsst/qserv-operator/api/v1alpha1"
	"github.com/lsst/qserv-operator/pkg/constants"
	"github.com/lsst/qserv-operator/pkg/syncer"
	"github.com/lsst/qserv-operator/pkg/syncers"
	"github.com/lsst/qserv-operator/pkg/util"
)

// QservReconciler reconciles a Qserv object
type QservReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps,resources=statefulsets;deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get;list;create;update
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=services;services/finalizers;configmaps;secrets,verbs=create;delete;get;list;patch;update;watch
// +kubebuilder:rbac:groups=networking.k8s.io,resources=networkpolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=qserv.lsst.org,resources=qservs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=qserv.lsst.org,resources=qservs/status,verbs=get;update;patch

// Reconcile reconciles a Qserv object
func (r *QservReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	log := r.Log.WithValues("qserv", req.NamespacedName)

	log.Info("Reconciling Qserv")

	// Fetch the Qserv instance
	qserv := &qservv1alpha1.Qserv{}
	err := r.Get(ctx, req.NamespacedName, qserv)
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

	result, err := r.updateQservStatus(ctx, req, qserv, &log)
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

// SetupWithManager setups Qserv controller for k8s
func (r *QservReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&qservv1alpha1.Qserv{}).
		Owns(&appsv1.StatefulSet{}).
		Complete(r)
}

func (r *QservReconciler) sync(syncers []syncer.Interface) error {
	for _, s := range syncers {
		if err := syncer.Sync(context.TODO(), s); err != nil {
			return err
		}
	}
	return nil
}

func (r *QservReconciler) updateQservStatus(ctx context.Context, req ctrl.Request, qserv *qservv1alpha1.Qserv, log *logr.Logger) (ctrl.Result, error) {
	// Manage status
	// See https://book.kubebuilder.io/cronjob-tutorial/controller-implementation.html#2-list-all-active-jobs-and-update-the-status
	listOpts := []client.ListOption{
		client.InNamespace(qserv.Namespace),
		client.MatchingLabels(util.GetInstanceLabels(qserv.Name)),
	}

	var statefulsets appsv1.StatefulSetList
	if err := r.List(ctx, &statefulsets, listOpts...); err != nil {
		(*log).Error(err, "Unable to list Qserv statefulsets")
		return ctrl.Result{}, err
	}
	hasStatefulSet := false
	hasDeployment := false
	var notReadyStatefulSet []appsv1.StatefulSet
	for _, statefulset := range statefulsets.Items {
		hasStatefulSet = true
		readyReplicas := statefulset.Status.ReadyReplicas
		desiredReplicas := *statefulset.Spec.Replicas
		readyFraction := fmt.Sprintf("%d/%d", readyReplicas, desiredReplicas)
		(*log).Info(fmt.Sprintf("Statefulset: %v, %s", statefulset.Name, readyFraction))
		if readyReplicas != desiredReplicas {
			notReadyStatefulSet = append(notReadyStatefulSet, statefulset)
		}

		componentLabel := statefulset.Labels["component"]
		switch componentLabel {
		case string(constants.Czar):
			qserv.Status.CzarReadyFraction = readyFraction
		case string(constants.IngestDb):
			qserv.Status.IngestDatabaseReadyFraction = readyFraction
		case string(constants.ReplCtl):
			qserv.Status.ReplicationControllerReadyFraction = readyFraction
		case string(constants.ReplDb):
			qserv.Status.ReplicationDatabaseReadyFraction = readyFraction
		case string(constants.Worker):
			qserv.Status.WorkerReadyFraction = readyFraction
		case string(constants.XrootdRedirector):
			qserv.Status.XrootdReadyFraction = readyFraction
		default:
			(*log).Info(fmt.Sprintf("Statefulset: %s has unknown 'component' label", statefulset.Name))
		}
	}

	// For qserv-dashboard deployment
	var deployments appsv1.DeploymentList
	if err := r.List(ctx, &deployments, listOpts...); err != nil {
		(*log).Error(err, "Unable to list Qserv deployments")
		return ctrl.Result{}, err
	}
	var notAvailableDeployment []appsv1.Deployment
	for _, deployment := range deployments.Items {
		hasDeployment = true
		availableReplicas := deployment.Status.AvailableReplicas
		desiredReplicas := *deployment.Spec.Replicas
		(*log).Info(fmt.Sprintf("Deployment: %v, %d/%d\n", deployment.Name, availableReplicas, desiredReplicas))
		if availableReplicas != desiredReplicas {
			notAvailableDeployment = append(notAvailableDeployment, deployment)
		}
	}

	availableCondition := metav1.Condition{
		Status: metav1.ConditionUnknown,
		Type:   "Available",
		Reason: "Succeed",
	}
	if !hasStatefulSet || !hasDeployment {
		availableCondition.Status = metav1.ConditionFalse
		availableCondition.Reason = "NotCreatedObjects"
		availableCondition.Message = "Statefulsets and deployment not yet created"
	} else if len(notReadyStatefulSet) != 0 || len(notAvailableDeployment) != 0 {
		availableCondition.Status = metav1.ConditionFalse
		availableCondition.Reason = "NotReadyPods"
		availableCondition.Message = "Pod(s) not ready or not available"
	} else {
		availableCondition.Status = metav1.ConditionTrue
	}
	availableCondition.LastTransitionTime = metav1.Now()

	qserv.Status.Conditions = []metav1.Condition{availableCondition}
	(*log).Info(fmt.Sprintf("Update status %v", qserv.Status.Conditions))
	err := r.Status().Update(ctx, qserv)
	return ctrl.Result{}, err
}
