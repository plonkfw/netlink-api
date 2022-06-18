package addrv1

import (
	"fmt"
	"net/http"

	utilsv1 "github.com/plonkfw/netlink-api/utils/v1"
	"github.com/vishvananda/netlink"
)

// ByName retrieves a address list by link name
func ByName(w http.ResponseWriter, r *http.Request) {
	// Get the name of the link from query params - /v1/addr/by/name?name=br0
	name := r.URL.Query().Get("name")

	// Lookup the link by name
	link, err := netlink.LinkByName(name)
	if err != nil {
		msg := fmt.Sprintf("Error querying link %s", name)
		utilsv1.Log.Error().Err(err).Msg(msg)
		utilsv1.ReplyError(w, r, msg, "ELOOKUPFAIL", err)
		return
	}

	// Get address info
	addressList, err := netlink.AddrList(link, 0)
	if err != nil {
		msg := fmt.Sprintf("Error refreshing info for link %s", name)
		utilsv1.Log.Error().Err(err).Msg(msg)
		utilsv1.ReplyError(w, r, msg, "ELOOKUPFAIL", err)
		return
	}

	// Prep response
	msg := fmt.Sprintf("Found addresses on link %s", name)
	utilsv1.ReplySuccess(w, r, msg, addressList)
}
