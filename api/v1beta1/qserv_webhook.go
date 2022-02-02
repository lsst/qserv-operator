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

package v1beta1

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	apivalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var log = logf.Log.WithName("webhook").WithName("qserv")

func (r *Qserv) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-qserv-lsst-org-v1beta1-qserv,mutating=true,failurePolicy=fail,sideEffects=None,groups=qserv.lsst.org,resources=qservs,verbs=create;update,versions=v1beta1,name=mqserv.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Qserv{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Qserv) Default() {
	log.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-qserv-lsst-org-v1beta1-qserv,mutating=false,failurePolicy=fail,sideEffects=None,groups=qserv.lsst.org,resources=qservs,verbs=create;update,versions=v1beta1,name=vqserv.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Qserv{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Qserv) ValidateCreate() error {
	log.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Qserv) ValidateUpdate(old runtime.Object) error {

	log.Info("validate update", "name", r.Name)
	// Validation logic upon object update.
	oldQservSpec, _ := old.(*Qserv)
	if err := r.validateQservUpdate(*oldQservSpec); err != nil {
		return err
	}

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Qserv) ValidateDelete() error {
	log.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

func (r *Qserv) validateQservUpdate(old Qserv) error {
	fldPath := field.NewPath("metadata")
	allErrs := apivalidation.ValidateObjectMetaUpdate(&r.ObjectMeta, &old.ObjectMeta, fldPath)
	qservSpec := old.Spec
	specErrs := r.validateQservSpecUpdate(qservSpec)
	allErrs = append(allErrs, specErrs...)

	if len(allErrs) == 0 {
		return nil
	}
	return apierrors.NewInvalid(
		schema.GroupKind{Group: "qserv.lsst.org", Kind: "Qserv"},
		r.Name, allErrs)
}

func (r *Qserv) validateQservSpecUpdate(oldQservSpec QservSpec) field.ErrorList {
	specErrs := field.ErrorList{}
	// Allow only additions to tolerations updates.
	// allErrs = append(allErrs, apivalidation.validateOnlyAddedTolerations(newPod.Spec.Tolerations, oldPod.Spec.Tolerations, specPath.Child("tolerations"))...)

	mungedPodSpec := *r.Spec.DeepCopy()
	// tolerations are checked before the deep copy, so munge those too
	mungedPodSpec.Image = oldQservSpec.Image
	mungedPodSpec.DbImage = oldQservSpec.DbImage
	mungedPodSpec.Tolerations = oldQservSpec.Tolerations // +k8s:verify-mutation:reason=clone
	log.Info("validate ", "name", r.Name)
	fieldPath := field.NewPath("spec")
	if !apiequality.Semantic.DeepEqual(oldQservSpec, mungedPodSpec) {
		if mungedPodSpec.Xrootd.Replicas != oldQservSpec.Xrootd.Replicas {
			fieldPath = field.NewPath("spec").Child("xrootd").Child("replicas")
		}
		// This diff isn't perfect, but it's a helluva lot better an "I'm not going to tell you what the difference is".
		//TODO: Pinpoint the specific field that causes the invalid error after we have strategic merge diff
		specDiff := cmp.Diff(mungedPodSpec, oldQservSpec)
		err := field.Forbidden(fieldPath, fmt.Sprintf("Qserv updates may not change fields other than "+
			"`spec.czar.image`, `spec.ingest.image`, `spec.replication.image`, `spec.worker.image`, `spec.xrootd.image`"+
			"`spec.tolerations` (only additions to existing tolerations)\n%v", specDiff))
		specErrs = append(specErrs, err)
	}
	return specErrs
}
