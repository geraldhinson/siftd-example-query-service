package routers

import (
	"github.com/geraldhinson/siftd-base/pkg/security"
	"github.com/geraldhinson/siftd-base/pkg/serviceBase"
	"github.com/geraldhinson/siftd-queryservice-base/pkg/queryhelpers"
)

type SecuredQueriesRouter struct {
	*queryhelpers.SecuredQueriesRoutesHelper
}

func NewSecuredQueriesRouter(queryService *serviceBase.ServiceBase) *SecuredQueriesRouter {

	// This router is mostly built using serviceBase implementation, but we don't fully implement it in
	// serviceBase because:
	// 1. We want it to be obvious that it is one of the routers this service implements (ie.
	//    easily seen in the routers folder)
	// 2. It is important for the service writer to define the auth model for all routers
	//

	//---------------------------------------------
	// queries that require an identity
	//
	// identity required queries can be called with MATCHING_IDENTITY or MACHINE_IDENTITY
	//
	var authModelUsers = queryService.NewAuthModel(security.ONE_DAY, []security.AuthTypes{security.MACHINE_IDENTITY, security.MATCHING_IDENTITY}, nil)
	if authModelUsers == nil {
		queryService.Logger.Fatalf("Failed to initialize AuthModelUsers")
	}

	//---------------------------------------------
	// queries that don't include an identity
	//
	// queries that lack an identity can be called with MACHINE_IDENTITY only
	//
	var authModelMachine = queryService.NewAuthModel(security.ONE_DAY, []security.AuthTypes{security.MACHINE_IDENTITY}, nil)
	if authModelMachine == nil {
		queryService.Logger.Fatalf("Failed to initialize AuthModelMachine")
	}

	// OK. Auth is defined. Now use the helper code to do the rest of the heavy lifting here.
	//
	SecuredQueriesRoutesHelper := queryhelpers.NewSecuredQueriesRoutesHelper(queryService, authModelUsers, authModelMachine)
	if SecuredQueriesRoutesHelper == nil {
		queryService.Logger.Println("Error creating PrivateQueriesRoutesHelper")
		return nil
	}

	SecuredQueriesRouter := &SecuredQueriesRouter{
		SecuredQueriesRoutesHelper: SecuredQueriesRoutesHelper,
	}

	return SecuredQueriesRouter
}
