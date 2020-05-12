package qserv

import (
	"context"

	qservv1alpha1 "github.com/lsst/qserv-operator/pkg/apis/qserv/v1alpha1"

	"github.com/lsst/qserv-operator/pkg/constants"
	"github.com/lsst/qserv-operator/pkg/controller/qserv/internal/sync"
	"github.com/lsst/qserv-operator/pkg/staging/syncer"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	kubedbscheme "kubedb.dev/apimachinery/client/clientset/versioned/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_qserv")

// Add creates a new Qserv Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileQserv{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("qserv-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Qserv
	err = c.Watch(&source.Kind{Type: &qservv1alpha1.Qserv{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Qserv
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &qservv1alpha1.Qserv{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileQserv implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileQserv{}

// ReconcileQserv reconciles a Qserv object
type ReconcileQserv struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client   client.Client
	scheme   *runtime.Scheme
	recorder record.EventRecorder
	// failover *failover.QservFailover
}

// Reconcile reads that state of the cluster for a Qserv object and makes changes based on the state read
// and what is in the Qserv.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileQserv) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Qserv")

	// Fetch the Qserv instance
	qserv := &qservv1alpha1.Qserv{}
	err := r.client.Get(context.TODO(), request.NamespacedName, qserv)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	r.scheme.Default(qserv)
	kubedbscheme.AddToScheme(r.scheme)

	syncers := []syncer.Interface{
		sync.NewCzarStatefulSetSyncer(qserv, r.client, r.scheme),
		sync.NewDotQservConfigMapSyncer(qserv, r.client, r.scheme),
		sync.NewWorkerStatefulSetSyncer(qserv, r.client, r.scheme),
		sync.NewReplicationCtlServiceSyncer(qserv, r.client, r.scheme),
		sync.NewReplicationCtlStatefulSetSyncer(qserv, r.client, r.scheme),
		sync.NewReplicationDbServiceSyncer(qserv, r.client, r.scheme),
		sync.NewReplicationDbStatefulSetSyncer(qserv, r.client, r.scheme),
		sync.NewXrootdRedirectorServiceSyncer(qserv, r.client, r.scheme),
		sync.NewXrootdStatefulSetSyncer(qserv, r.client, r.scheme),
	}

	syncers = append(sync.NewQservServicesSyncer(qserv, r.client, r.scheme), syncers...)

	// Redis database: optional, stores secondary index data
	if qserv.Spec.Redis != nil {
		reqLogger.Info("Reconciling Redis")
		syncers = append(syncers, sync.NewRedisSyncer(qserv, r.client, r.scheme))
	}

	for _, configmapClass := range constants.ContainerConfigmaps {
		for _, subpath := range []string{"etc", "start"} {
			syncers = append(syncers, sync.NewContainerConfigMapSyncer(qserv, r.client, r.scheme, configmapClass, subpath))
		}
	}
	syncers = append(syncers, sync.NewContainerConfigMapSyncer(qserv, r.client, r.scheme, constants.InitDbName, "start"))

	for _, db := range constants.Databases {
		syncers = append(syncers, sync.NewSqlConfigMapSyncer(qserv, r.client, r.scheme, db))
	}

	// Specify Network Policies
	syncers = append(syncers, sync.NewNetworkPoliciesSyncer(qserv, r.client, r.scheme)...)

	if err = r.sync(syncers); err != nil {
		return reconcile.Result{}, err
	}

	// if err = r.failover.CheckAndHeal(qserv); err != nil {
	// 	return reconcile.Result{}, err
	// }

	return reconcile.Result{}, nil
}

func (r *ReconcileQserv) sync(syncers []syncer.Interface) error {
	for _, s := range syncers {
		if err := syncer.Sync(context.TODO(), s, r.recorder); err != nil {
			return err
		}
	}
	return nil
}
