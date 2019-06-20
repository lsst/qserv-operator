package controller

import (
	"github.com/lsst/qserv-operator/pkg/controller/qserv"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, qserv.Add)
}
