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
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	qservv1beta1 "github.com/lsst/qserv-operator/api/v1beta1"
	"github.com/lsst/qserv-operator/controllers/constants"
	"github.com/lsst/qserv-operator/controllers/specs"
	"github.com/lsst/qserv-operator/controllers/util"
)

// QservReconciler reconciles a Qserv object
type QservReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="networking.k8s.io",resources=networkpolicies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
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
	log.V(5).Enabled()
	// TODO check which log to use
	// log := r.Log.WithValues("qserv", request.NamespacedName)

	log.V(0).Info("Reconcile Qserv")

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

	result, err := r.updateQservStatus(ctx, request, qserv)
	if err != nil {
		log.Error(err, "Unable to update Qserv status")
		return result, err
	}

	objectSpecManagers := []ObjectSpecManager{
		&specs.CzarSpec{},
		&specs.CzarServiceSpec{},
		&specs.DotQservConfigMapSpec{},
		&specs.IngestDatabaseSpec{},
		&specs.IngestDatabaseServiceSpec{},
		&specs.QueryServiceSpec{},
		&specs.ReplicationControllerServiceSpec{},
		&specs.ReplicationControllerSpec{},
		&specs.ReplicationRegistrySpec{},
		&specs.ReplicationRegistryServiceSpec{},
		&specs.ReplicationDatabaseSpec{},
		&specs.ReplicationDatabaseServiceSpec{},
		&specs.WorkerServiceSpec{},
		&specs.WorkerSpec{},
		&specs.XrootdServiceSpec{},
		&specs.XrootdSpec{},
	}

	// Manage "*-etc" and "*-start" configmaps
	var configmapSpec ObjectSpecManager
	configmapSpec = &specs.ContainerConfigMapSpec{
		ContainerName: constants.InitDbName,
		Subdir:        "start",
	}

	objectSpecManagers = append(objectSpecManagers, configmapSpec)
	for _, containerName := range constants.WithEtcStartConfigmaps {
		for _, subdir := range []string{"etc", "start"} {
			configmapSpec = &specs.ContainerConfigMapSpec{
				ContainerName: containerName,
				Subdir:        subdir,
			}
			objectSpecManagers = append(objectSpecManagers, configmapSpec)
		}
	}

	for _, containerName := range constants.WithStartConfigmap {
		configmapSpec = &specs.ContainerConfigMapSpec{
			ContainerName: containerName,
			Subdir:        "start",
		}
		objectSpecManagers = append(objectSpecManagers, configmapSpec)
	}

	for _, database := range constants.Databases {
		configmapSpec = &specs.SQLConfigMapSpec{
			Database: database,
		}
		objectSpecManagers = append(objectSpecManagers, configmapSpec)
	}

	// Create Network Policies specification
	if qserv.Spec.NetworkPolicies {
		networkPolicySpecManagers := []ObjectSpecManager{
			&specs.CzarNetworkPolicySpec{},
			&specs.DefaultNetworkPolicySpec{},
			&specs.ReplDatabaseNetworkPolicySpec{},
			&specs.WorkerNetworkPolicySpec{},
			&specs.XrootdRedirectorNetworkPolicySpec{},
		}
		objectSpecManagers = append(objectSpecManagers, networkPolicySpecManagers...)
	}

	// Reconcile all objects
	for _, objectSpecManager := range objectSpecManagers {
		result, err = r.reconcile(ctx, qserv, objectSpecManager)
		if err != nil {
			log.Error(err, "Unable to reconcile", "name", objectSpecManager.GetName())
			return result, err
		}
	}

	// TODO: understand event management and implement it
	// see: https://github.com/kubernetes-sigs/kubebuilder/discussions/2465

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
		Owns(&v1.ConfigMap{}).
		Owns(&v1.Service{}).
		Owns(&appsv1.Deployment{}).
		Owns(&appsv1.StatefulSet{}).
		Complete(r)
}

func (r *QservReconciler) updateQservStatus(ctx context.Context, req ctrl.Request, qserv *qservv1beta1.Qserv) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	// Manage status
	// See https://book.kubebuilder.io/cronjob-tutorial/controller-implementation.html#2-list-all-active-jobs-and-update-the-status
	listOpts := []client.ListOption{
		client.InNamespace(qserv.Namespace),
		client.MatchingLabels(util.GetInstanceLabels(qserv.Name)),
	}

	var statefulsets appsv1.StatefulSetList
	if err := r.List(ctx, &statefulsets, listOpts...); err != nil {
		log.Error(err, "Unable to list Qserv StatefulSets")
		return ctrl.Result{}, err
	}
	hasStatefulSet := false
	var notReadyStatefulSet []appsv1.StatefulSet
	for _, statefulset := range statefulsets.Items {
		hasStatefulSet = true
		readyReplicas := statefulset.Status.ReadyReplicas
		desiredReplicas := *statefulset.Spec.Replicas
		readyFraction := fmt.Sprintf("%d/%d", readyReplicas, desiredReplicas)
		log.Info("Check resource status", "resource kind", "Statefulset", "resource name", statefulset.Name, "ready fraction", readyFraction)
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
			log.Info("Unknown label", "label-name", "component", "kind", "statefulset", "name", statefulset.Name, "ready fraction", readyFraction)
		}
	}

	var deployments appsv1.DeploymentList
	if err := r.List(ctx, &deployments, listOpts...); err != nil {
		log.Error(err, "Unable to list Qserv Deployments")
		return ctrl.Result{}, err
	}
	hasDeployment := false
	var notReadyDeployment []appsv1.Deployment
	for _, deployment := range deployments.Items {
		hasDeployment = true
		readyReplicas := deployment.Status.ReadyReplicas
		desiredReplicas := *deployment.Spec.Replicas
		readyFraction := fmt.Sprintf("%d/%d", readyReplicas, desiredReplicas)
		log.Info("Check resource status", "resource kind", "Deployment", "resource name", deployment.Name, "ready fraction", readyFraction)
		if readyReplicas != desiredReplicas {
			notReadyDeployment = append(notReadyDeployment, deployment)
		}

		componentLabel := deployment.Labels["component"]
		switch componentLabel {
		case string(constants.ReplRegistry):
			qserv.Status.ReplicationRegistryReadyFraction = readyFraction
		default:
			log.Info("Unknown label", "label-name", "component", "kind", "deployment", "name", deployment.Name, "ready fraction", readyFraction)
		}
	}

	availableCondition := metav1.Condition{
		Status: metav1.ConditionUnknown,
		Type:   "Available",
		Reason: "Succeed",
	}
	if !hasStatefulSet {
		availableCondition.Status = metav1.ConditionFalse
		availableCondition.Reason = "NotCreatedStatefulsets"
		availableCondition.Message = "Statefulsets not yet created"
	} else if !hasDeployment {
		availableCondition.Status = metav1.ConditionFalse
		availableCondition.Reason = "NotCreatedDeployment"
		availableCondition.Message = "Deployment not yet created"
	} else if len(notReadyStatefulSet) != 0 || len(notReadyDeployment) != 0 {
		availableCondition.Status = metav1.ConditionFalse
		availableCondition.Reason = "NotReadyPods"
		availableCondition.Message = "Pod(s) not ready or not available"
	} else {
		availableCondition.Status = metav1.ConditionTrue
	}
	availableCondition.LastTransitionTime = metav1.Now()

	qserv.Status.Conditions = []metav1.Condition{availableCondition}
	log.Info("Update Qserv conditions", "conditions", qserv.Status.Conditions)
	err := r.Status().Update(ctx, qserv)
	return ctrl.Result{}, err
}
