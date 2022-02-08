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
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	qservv1beta1 "github.com/lsst/qserv-operator/api/v1beta1"
)

// ObjectSpecManager manage creation and update of specifications for standard k8s API objects
// related to Qserv
type ObjectSpecManager interface {
	Create() (client.Object, error)
	GetName() string
	Initialize(qserv *qservv1beta1.Qserv) client.Object
	Update(object client.Object) (bool, error)
}

func (r *QservReconciler) reconcile(ctx context.Context, qserv *qservv1beta1.Qserv, controlled ObjectSpecManager) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithName("reconcile")
	// Check if the current controlled API object exists, if not create it
	object := controlled.Initialize(qserv)
	key := types.NamespacedName{Name: controlled.GetName(), Namespace: qserv.Namespace}
	err := r.Get(ctx, key, object)

	if err != nil {
		if errors.IsNotFound(err) {
			// Define and create a new object.
			object, err = controlled.Create()
			if err != nil {
				return ctrl.Result{}, err
			}
			if err = controllerutil.SetControllerReference(qserv, object, r.Scheme); err != nil {
				return ctrl.Result{}, err
			}
			log.V(0).Info("Create ", "key", key)
			if err = r.Create(ctx, object); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, err

	}

	// Check if the current controlled API object require an update, and then perform the update
	requireUpdate, err2 := controlled.Update(object)
	if err2 != nil {
		return ctrl.Result{}, err2
	}
	if requireUpdate {
		log.V(0).Info("Update ", "key", key)
		if err = r.Update(ctx, object); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	return ctrl.Result{}, err
}
