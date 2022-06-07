package link

import (
	"encoding/json"
	"net/http"

	"github.com/plonkfw/netlink-api/utils"
	"github.com/vishvananda/netlink"
)

// List returns a list of all interfaces
func List(w http.ResponseWriter, r *http.Request) {
	// Fetch a list of all interfaces
	links, err := netlink.LinkList()
	if err != nil {
		utils.Log.Error().Err(err).Msg("Error listing liks")
	}

	// Prepare the response
	response, err := json.MarshalIndent(links, "", "  ")
	if err != nil {
		utils.Log.Error().Err(err).Msg("Error marshling response")
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
