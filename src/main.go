package main

import (
	"net/http"
	"time"

	routingv1 "github.com/plonkfw/netlink-api/routing/v1"
	utilsv1 "github.com/plonkfw/netlink-api/utils/v1"
)

// Start logging
func init() {
	t := time.Now()
	// WHY isn't this format specification just yyyy-MM-dd etc
	utilsv1.Log.Info().Msg("Logs begin at " + t.UTC().Format("2006-01-02 15:04:05") + " UTC")
}

// Setup http listener and routes
func main() {
	routerv1 := routingv1.NewAPIRouter()
	utilsv1.Log.Fatal().Err(http.ListenAndServe(":4821", routerv1))
}
