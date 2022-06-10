package linkv1

import (
	"fmt"
	"net/http"

	utilsv1 "github.com/plonkfw/netlink-api/utils/v1"
	"github.com/vishvananda/netlink"
)

// List returns a list of all interfaces
func List(w http.ResponseWriter, r *http.Request) {
	// Fetch a list of all interfaces
	links, err := netlink.LinkList()
	if err != nil {
		utilsv1.Log.Error().Err(err).Msg("Error listing links")
	}

	msg := fmt.Sprintf("Found interfaces")
	utilsv1.ReplySuccess(w, r, msg, links)
}
