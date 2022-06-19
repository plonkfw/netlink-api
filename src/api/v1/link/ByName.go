package linkv1

import (
	"errors"
	"fmt"
	"net/http"

	utilsv1 "github.com/plonkfw/netlink-api/utils/v1"
	"github.com/vishvananda/netlink"
)

// ByName retrieves a link by name
func ByName(w http.ResponseWriter, r *http.Request) {
	// Get the name of the link from query params - /v1/link/by-name?name=lo
	name := r.URL.Query().Get("name")

	// Did they provide a param
	if name != "" {
		// Lookup the link by name
		link, err := netlink.LinkByName(name)
		if err != nil {
			msg := fmt.Sprintf("Error querying link %s", name)
			utilsv1.Log.Error().Err(err).Msg(msg)
			utilsv1.ReplyError(w, r, msg, "ELOOKUPFAIL", err)
			return
		}

		// Prep response
		msg := fmt.Sprintf("Found interface %s", name)
		utilsv1.ReplySuccess(w, r, msg, link)
	}

	// Invalid params
	msg := fmt.Sprintf("Invalid paramaters %s", name)
	err := errors.New(msg)
	utilsv1.Log.Error().Err(err).Msg(msg)
	utilsv1.ReplyError(w, r, msg, "EINVALIDPARAM", err)
	return
}
