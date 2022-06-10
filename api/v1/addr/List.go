package addrv1

import (
	"fmt"
	"net/http"

	utilsv1 "github.com/plonkfw/netlink-api/utils/v1"
	"github.com/vishvananda/netlink"
)

// List returns a list of all addresses in the system
func List(w http.ResponseWriter, r *http.Request) {
	// Fetch a list of all addresses
	addresses, err := netlink.AddrList(nil, 0)
	if err != nil {
		utilsv1.Log.Error().Err(err).Msg("Error listing addresses")
	}

	msg := fmt.Sprintf("Found addresses")
	utilsv1.ReplySuccess(w, r, msg, addresses)
}
