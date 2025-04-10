package routers

import (
	"github.com/geraldhinson/siftd-base/pkg/security"
	"github.com/geraldhinson/siftd-base/pkg/serviceBase"
	"github.com/geraldhinson/siftd-queryservice-base/pkg/queryhelpers"
)

type HealthCheckRouter struct {
	*queryhelpers.HealthCheckRoutesHelper
}

func NewHealthCheckRouter(queryService *serviceBase.ServiceBase) *HealthCheckRouter {

	// This router is mostly built using query-servvice implementation, but we don't fully implement it in
	// query-service because:
	// 1. We want it to be obvious that it is one of the routers this service implements (ie.
	//    easily seen in the routers folder)
	// 2. It is important for the service writer to define the auth model for all routers
	//
	authModel, err := queryService.NewAuthModel(security.NO_REALM, security.NO_AUTH, security.NO_EXPIRY, nil)
	if err != nil {
		queryService.Logger.Fatalf("Failed to initialize AuthModelUsers in HealthCheckRouter : %v", err)
		return nil
	}

	// OK. Auth is defined. Now use the helper code to do the rest of the heavy lifting here.
	//
	HealthCheckRoutesHelper := queryhelpers.NewHealthCheckRoutesHelper(queryService, authModel)
	if HealthCheckRoutesHelper == nil {
		queryService.Logger.Println("Error creating HealthCheckRoutesHelper")
		return nil
	}

	HealthCheckServiceRouter := &HealthCheckRouter{
		HealthCheckRoutesHelper: HealthCheckRoutesHelper,
	}

	return HealthCheckServiceRouter
}
