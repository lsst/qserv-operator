package util

import (
	"fmt"

	qservv1alpha1 "github.com/lsst/qserv-operator/pkg/apis/qserv/v1alpha1"
)

// GetName returns a name whose prefix is instance name and suffix typeName
func GetName(r *qservv1alpha1.Qserv, typeName string) string {
	return fmt.Sprintf("%s-%s", r.Name, typeName)
}
