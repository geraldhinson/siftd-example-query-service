package routers

import (
	"github.com/geraldhinson/siftd-base/pkg/helpers"
	"github.com/geraldhinson/siftd-base/pkg/security"
	"github.com/geraldhinson/siftd-base/pkg/serviceBase"
)

type FakeIdentityServiceRouter struct {
	*helpers.FakeIdentityServiceRoutesHelper
}

func NewFakeIdentityServiceRouter(queryService *serviceBase.ServiceBase) *FakeIdentityServiceRouter {

	// This router is mostly built using serviceBase implementation, but we don't fully implement it in
	// serviceBase because:
	// 1. We want it to be obvious that it is one of the routers this service implements (ie.
	//    easily seen in the routers folder)
	// 2. It is important for the service writer to define the auth model for all routers
	//
	authModel, err := queryService.NewAuthModel(security.NO_REALM, security.NO_AUTH, security.NO_EXPIRY, nil)
	if err != nil {
		queryService.Logger.Fatalf("Failed to initialize AuthModelUsers in FakeIdentityServiceRouter : %v", err)
		return nil
	}

	// OK. Auth is defined. Now use the helper code to do the rest of the heavy lifting here.
	//
	FakeIdentityServiceRoutesHelper := helpers.NewFakeIdentityServiceRoutesHelper(queryService, authModel)
	if FakeIdentityServiceRoutesHelper == nil {
		queryService.Logger.Println("Error creating FakeIdentityServiceRoutesHelper")
		return nil
	}

	FakeIdentityServiceRouter := &FakeIdentityServiceRouter{
		FakeIdentityServiceRoutesHelper: FakeIdentityServiceRoutesHelper,
	}

	return FakeIdentityServiceRouter
}
