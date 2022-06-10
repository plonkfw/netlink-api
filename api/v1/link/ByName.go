package linkv1

import (
	"fmt"
	"net/http"

	"github.com/plonkfw/netlink-api/utils"
	"github.com/vishvananda/netlink"
)

// ByName retrieves a link by name
func ByName(w http.ResponseWriter, r *http.Request) {
	// Get the name of the link from query params - /v1/link/by-name?name=lo
	name := r.URL.Query().Get("name")

	// Lookup the link by name
	link, err := netlink.LinkByName(name)
	if err != nil {
		msg := fmt.Sprintf("Error querying link %s", name)
		utils.Log.Error().Err(err).Msg(msg)
		utils.ReplyError(w, r, msg, err)
		return
	}

	// Prep response
	msg := fmt.Sprintf("Found interface %s", name)
	utils.ReplySuccess(w, r, msg, link)
}
