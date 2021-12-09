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

	"github.com/go-logr/logr"
	qservv1beta1 "github.com/lsst/qserv-operator/api/v1beta1"
	"github.com/lsst/qserv-operator/controllers/constants"
)

type ObjectSpecManager interface {
	CreateObjectSpec() client.Object
	Create(qserv *qservv1beta1.Qserv, object *client.Object) error
	Update(qserv *qservv1beta1.Qserv, object *client.Object) (bool, error)
}

func (r *QservReconciler) reconcile(ctx context.Context, qserv *qservv1beta1.Qserv, log logr.Logger, controlled ObjectSpecManager) (ctrl.Result, error) {
	// Check if the czar statefulset already exists, if not create a new statefulset.
	object := controlled.CreateObjectSpec()
	err := r.Get(ctx, types.NamespacedName{Name: qserv.Name + "-" + string(constants.Czar), Namespace: qserv.Namespace}, object)

	if err != nil {
		if errors.IsNotFound(err) {
			// Define and create a new object.
			if err = controlled.Create(qserv, &object); err != nil {
				return ctrl.Result{}, err
			}
			controllerutil.SetControllerReference(qserv, object, r.Scheme)
			log.V(0).Info("Create Qserv")
			if err = r.Create(ctx, object); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil
		} else {
			return ctrl.Result{}, err
		}
	}

	// Ensure the statefulset size is the same as the spec.
	update, err2 := controlled.Update(qserv, &object)
	if err2 != nil {
		return ctrl.Result{}, err2
	}
	if update {
		log.V(0).Info("Update Qserv")
		if err = r.Update(ctx, object); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	return ctrl.Result{}, err
}
