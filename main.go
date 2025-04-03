package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/geraldhinson/siftd-base/pkg/constants"
	"github.com/geraldhinson/siftd-base/pkg/serviceBase"
	"github.com/geraldhinson/siftd-example-query-service/routers"
)

func main() {
	// call setup to get the service base (logging, config, routing) and a keystore
	queryService := serviceBase.NewServiceBase()
	if queryService == nil {
		fmt.Println("Failed to validate configuration and listen. Shutting down.")
		return
	}

	SecuredQueriesRouter := routers.NewSecuredQueriesRouter(queryService)
	if SecuredQueriesRouter == nil {
		queryService.Logger.Fatalf("Failed to create secured queries api server. Shutting down.")
		return
	}

	PublicQueriesRouter := routers.NewPublicQueriesRouter(queryService)
	if PublicQueriesRouter == nil {
		queryService.Logger.Fatalf("Failed to create public queries api server. Shutting down.")
		return
	}

	HealthCheckRouter := routers.NewHealthCheckRouter(queryService)
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
		FakeIdentityServiceRouter := routers.NewFakeIdentityServiceRouter(queryService)
		if FakeIdentityServiceRouter == nil {
			queryService.Logger.Fatalf("Failed to create fake identity service api server (for testing only). Shutting down.")
			return
		}
	}

	queryService.Logger.Printf("This service is listening on: %s", listenAddress)
	http.ListenAndServe(listenAddress, queryService.Router)

}
