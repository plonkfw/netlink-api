package linkv1

import (
	"fmt"
	"net/http"
	"strconv"

	utilsv1 "github.com/plonkfw/netlink-api/utils/v1"
	"github.com/vishvananda/netlink"
)

// ByIndex retrieves a link by index
func ByIndex(w http.ResponseWriter, r *http.Request) {
	// Get the Index of the link from query params - /v1/link/by-name?name=lo
	query := r.URL.Query().Get("index")
	index, _ := strconv.Atoi(query)

	// Lookup the link by Index
	link, err := netlink.LinkByIndex(index)
	if err != nil {
		msg := fmt.Sprintf("Error querying link %d", index)
		utilsv1.Log.Error().Err(err).Msg(msg)
		utilsv1.ReplyError(w, r, msg, "ELOOKUPFAIL", err)
		return
	}

	// Prep response
	msg := fmt.Sprintf("Found interface %d", index)
	utilsv1.ReplySuccess(w, r, msg, link)
}
