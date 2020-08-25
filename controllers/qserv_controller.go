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
	"reflect"

	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	qservv1alpha1 "github.com/lsst/qserv-operator/api/v1alpha1"
	"github.com/lsst/qserv-operator/pkg/constants"
	"github.com/lsst/qserv-operator/pkg/syncer"
	"github.com/lsst/qserv-operator/pkg/syncers"
)

// QservReconciler reconciles a Qserv object
type QservReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []v1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

// labelsForQserv returns the labels for selecting the resources
// belonging to the given qserv CR name.
func labelsForQserv(name string) map[string]string {
	return map[string]string{"app": "qserv", "qserv_cr": name}
}

// +kubebuilder:rbac:groups=qserv.lsst.org,resources=qservs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=qserv.lsst.org,resources=qservs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list
// +kubebuilder:rbac:groups=core,resources=services;services/finalizers;configmaps;secrets,verbs=create;delete;get;list;patch;update;watch

func (r *QservReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {

	ctx := context.Background()
	log := r.Log.WithValues("memcached", req.NamespacedName)

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

	r.Scheme.Default(qserv)

	syncers := []syncer.Interface{
		syncers.NewCzarStatefulSetSyncer(qserv, r.client, r.scheme),
		syncers.NewDotQservConfigMapSyncer(qserv, r.client, r.scheme),
		syncers.NewWorkerStatefulSetSyncer(qserv, r.client, r.scheme),
		syncers.NewReplicationCtlServiceSyncer(qserv, r.client, r.scheme),
		syncers.NewReplicationCtlStatefulSetSyncer(qserv, r.client, r.scheme),
		syncers.NewIngestDbServiceSyncer(qserv, r.client, r.scheme),
		syncers.NewIngestDbStatefulSetSyncer(qserv, r.client, r.scheme),
		syncers.NewReplicationDbServiceSyncer(qserv, r.client, r.scheme),
		syncers.NewReplicationDbStatefulSetSyncer(qserv, r.client, r.scheme),
		syncers.NewXrootdRedirectorServiceSyncer(qserv, r.client, r.scheme),
		syncers.NewXrootdStatefulSetSyncer(qserv, r.client, r.scheme),
	}

	syncers = append(syncers.NewQservServicesSyncer(qserv, r.client, r.scheme), syncers...)

	// Redis database: optional, stores secondary index data
	if qserv.Spec.Redis != nil {
		reqLogger.Info("Reconciling Redis")
		syncers = append(syncers, syncers.NewRedisSyncer(qserv, r.client, r.scheme))
	}

	for _, configmapClass := range constants.ContainerConfigmaps {
		for _, subpath := range []string{"etc", "start"} {
			syncers = append(syncers, syncers.NewContainerConfigMapSyncer(qserv, r.client, r.scheme, configmapClass, subpath))
		}
	}
	syncers = append(syncers, syncers.NewContainerConfigMapSyncer(qserv, r.client, r.scheme, constants.InitDbName, "start"))

	for _, db := range constants.Databases {
		syncers = append(syncers, syncers.NewSqlConfigMapSyncer(qserv, r.client, r.scheme, db))
	}

	// Specify Network Policies
	if qserv.Spec.NetworkPolicies {
		syncers = append(syncers, syncers.NewNetworkPoliciesSyncer(qserv, r.client, r.scheme)...)
	}

	if err = r.sync(syncers); err != nil {
		return reconcile.Result{}, err
	}

	// Update the Qserv status with the pod names
	// List the pods for this qserv's deployment
	podList := &v1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(qserv.Namespace),
		client.MatchingLabels(labelsForMemcached(qserv.Name)),
	}
	if err = r.List(ctx, podList, listOpts...); err != nil {
		log.Error(err, "Failed to list pods", "Qserv.Namespace", qserv.Namespace, "Qserv.Name", qserv.Name)
		return ctrl.Result{}, err
	}
	podNames := getPodNames(podList.Items)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(podNames, qserv.Status.Nodes) {
		qserv.Status.Nodes = podNames
		err := r.Status().Update(ctx, qserv)
		if err != nil {
			log.Error(err, "Failed to update Qserv status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *QservReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&qservv1alpha1.Qserv{}).
		Owns(&appsv1.Statefulset{}).
		Complete(r)
}

func (r *QservReconciler) sync(syncers []syncer.Interface) error {
	for _, s := range syncers {
		if err := syncer.Sync(context.TODO(), s, r.recorder); err != nil {
			return err
		}
	}
	return nil
}
