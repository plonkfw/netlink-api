package link

import (
	"fmt"
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

	msg := fmt.Sprintf("Found interfaces")
	utils.ReplySuccess(w, r, msg, links)
}
