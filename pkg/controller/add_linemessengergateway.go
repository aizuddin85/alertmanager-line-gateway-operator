package controller

import (
	"github.com/aizuddin85/alertmanager-line-gateway-operator/pkg/controller/linemessengergateway"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, linemessengergateway.Add)
}
