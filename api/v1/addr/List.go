package addr

import (
	"fmt"
	"net/http"

	"github.com/plonkfw/netlink-api/utils"
	"github.com/vishvananda/netlink"
)

// List returns a list of all addresses in the system
func List(w http.ResponseWriter, r *http.Request) {
	// Fetch a list of all addresses
	addresses, err := netlink.AddrList(nil, 0)
	if err != nil {
		utils.Log.Error().Err(err).Msg("Error listing addresses")
	}

	msg := fmt.Sprintf("Found addresses")
	utils.ReplySuccess(w, r, msg, addresses)
}
