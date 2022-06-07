package main

import (
	"net/http"
	"time"

	"github.com/plonkfw/netlink-api/routing"
	"github.com/plonkfw/netlink-api/utils"
)

// Start logging
func init() {
	t := time.Now()
	// WHY isn't this format specification just yyyy-MM-dd etc
	utils.Log.Info().Msg("Logs begin at " + t.UTC().Format("2006-01-02 15:04:05") + " UTC")
}

// Setup http listener and routes
func main() {
	router := routing.NewAPIRouter()
	utils.Log.Fatal().Err(http.ListenAndServe(":4821", router))
}
