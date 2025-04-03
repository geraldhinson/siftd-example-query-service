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
	var authModelUsers = queryService.NewAuthModel(security.ONE_DAY, []security.AuthTypes{security.NO_AUTH}, nil)
	if authModelUsers == nil {
		queryService.Logger.Fatalf("Failed to initialize AuthModelUsers in HealthCheckServiceRouter")
	}

	// OK. Auth is defined. Now use the helper code to do the rest of the heavy lifting here.
	//
	HealthCheckRoutesHelper := queryhelpers.NewHealthCheckRoutesHelper(queryService, authModelUsers)
	if HealthCheckRoutesHelper == nil {
		queryService.Logger.Println("Error creating HealthCheckRoutesHelper")
		return nil
	}

	HealthCheckServiceRouter := &HealthCheckRouter{
		HealthCheckRoutesHelper: HealthCheckRoutesHelper,
	}

	return HealthCheckServiceRouter
}
