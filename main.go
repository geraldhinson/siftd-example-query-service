package main

import (
	"fmt"
	"strings"

	"github.com/geraldhinson/siftd-base/pkg/constants"
	"github.com/geraldhinson/siftd-base/pkg/helpers"
	"github.com/geraldhinson/siftd-base/pkg/security"
	"github.com/geraldhinson/siftd-base/pkg/serviceBase"
	"github.com/geraldhinson/siftd-queryservice-base/pkg/queryhelpers"
)

func main() {
	// call setup to get the service base (logging, config, routing) and a keystore
	queryService := serviceBase.NewServiceBase()
	if queryService == nil {
		fmt.Println("Failed to validate configuration and listen. Shutting down.")
		return
	}

	SecuredQueriesRouter := queryhelpers.NewSecuredQueriesRouter(
		queryService,
		security.REALM_MEMBER, security.MATCHING_IDENTITY, security.ONE_DAY, nil, // queries that include an identity in the URL
		security.REALM_MACHINE, security.VALID_IDENTITY, security.ONE_HOUR, nil) // queries that don't include an identity in the URL
	if SecuredQueriesRouter == nil {
		queryService.Logger.Fatalf("Failed to create secured queries api server. Shutting down.")
		return
	}

	PublicQueriesRouter := queryhelpers.NewPublicQueriesRouter(queryService, security.NO_REALM, security.NO_AUTH, security.NO_EXPIRY, nil)
	if PublicQueriesRouter == nil {
		queryService.Logger.Fatalf("Failed to create public queries api server. Shutting down.")
		return
	}

	HealthCheckRouter := queryhelpers.NewHealthCheckRouter(queryService, security.NO_REALM, security.NO_AUTH, security.NO_EXPIRY, nil)
	if HealthCheckRouter == nil {
		queryService.Logger.Fatalf("Failed to create health check api server. Shutting down.")
		return
	}

	// setup
	listenAddress := queryService.Configuration.GetString(constants.LISTEN_ADDRESS)
	if listenAddress == "" {
		queryService.Logger.Fatalf("Unable to retrieve listen address and port. Shutting down.")
		return
	}

	// if we are running on localhost, we can add a fake identity service for testing (id is hardcoded in FakeKeyStore.go)
	if strings.Contains(listenAddress, "localhost") {
		FakeIdentityServiceRouter := helpers.NewFakeIdentityServiceRouter(queryService, security.NO_REALM, security.NO_AUTH, security.NO_EXPIRY, nil)
		if FakeIdentityServiceRouter == nil {
			queryService.Logger.Fatalf("Failed to create fake identity service api server (for testing only). Shutting down.")
			return
		}
	}

	queryService.ListenAndServe()
}
