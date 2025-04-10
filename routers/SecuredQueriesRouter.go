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
	authModelUsers, err := queryService.NewAuthModel(security.REALM_MEMBER, security.MATCHING_IDENTITY, security.ONE_DAY, nil)
	if err != nil {
		queryService.Logger.Fatalf("Failed to initialize AuthModelUsers in SecuredQueriesRouter : %v", err)
		return nil
	}

	//---------------------------------------------
	// queries that don't include an identity
	//
	// queries that lack an identity can be called with MACHINE_IDENTITY only
	//
	authModelMachine, err := queryService.NewAuthModel(security.REALM_MACHINE, security.VALID_IDENTITY, security.ONE_HOUR, nil)
	if err != nil {
		queryService.Logger.Fatalf("Failed to initialize AuthModelMachine in SecuredQueriesRouter : %v", err)
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
