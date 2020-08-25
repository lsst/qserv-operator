package syncer

import (
	"context"
	"fmt"

	"github.com/iancoleman/strcase"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("syncer")

const (
	eventNormal  = "Normal"
	eventWarning = "Warning"
)

func getKey(obj runtime.Object) (types.NamespacedName, error) {
	key := types.NamespacedName{}
	objMeta, ok := obj.(metav1.Object)
	if !ok {
		return key, fmt.Errorf("%T is not a metav1.Object", obj)
	}

	key.Name = objMeta.GetName()
	key.Namespace = objMeta.GetNamespace()
	return key, nil
}

func basicEventReason(objKindName string, err error) string {
	if err != nil {
		return fmt.Sprintf("%sSyncFailed", strcase.ToCamel(objKindName))
	}
	return fmt.Sprintf("%sSyncSuccessfull", strcase.ToCamel(objKindName))
}

// Sync mutates the subject of the syncer interface using controller-runtime
// CreateOrUpdate method, when obj is not nil. It takes care of setting owner
// references and recording kubernetes events where appropriate
func Sync(ctx context.Context, syncer Interface) error {
	_, err := syncer.Sync(ctx)
	return err
}

// WithoutOwner partially implements implements the syncer interface for the case the subject has no owner
type WithoutOwner struct{}

// GetOwner implementation of syncer interface for the case the subject has no owner
func (*WithoutOwner) GetOwner() runtime.Object {
	return nil
}
