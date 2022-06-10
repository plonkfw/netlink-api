package link

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/plonkfw/netlink-api/utils"
	"github.com/vishvananda/netlink"
)

type setNoMaster struct {
	Link string
}

// SetNoMaster removes the master of the link device
func SetNoMaster(w http.ResponseWriter, r *http.Request) {
	var setNoMaster setNoMaster

	// Unpack the request
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		utils.Log.Error().Err(err).Msg("Error reading body")
		utils.ReplyError(w, r, "Error reading body", err)
		return
	}

	// Print the request to deubg stream
	utils.Log.Debug().Msg(string(body))

	// Unpack the request
	if err := json.Unmarshal(body, &setNoMaster); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			utils.Log.Error().Err(err).Msg("Error unmarshaling body")
			utils.ReplyError(w, r, "Error unmarshaling body", err)
			return
		}
	}

	if setNoMaster.Link != "" {
		link, _ := netlink.LinkByName(setNoMaster.Link)
		if err != nil {
			msg := fmt.Sprintf("Error looking up link %s", setNoMaster.Link)
			utils.Log.Error().Err(err).Msg(msg)
			utils.ReplyError(w, r, msg, err)
			return
		}

		err = nil
		err = netlink.LinkSetNoMaster(link)
		if err != nil {
			msg := fmt.Sprintf("Error removing master of link %s", setNoMaster.Link)
			utils.Log.Error().Err(err).Msg(msg)
			utils.ReplyError(w, r, msg, err)
			return
		}

		// Lookup the link by name
		refreshedLink, _ := netlink.LinkByName(setNoMaster.Link)

		// Prep response
		msg := fmt.Sprintf("Successfully removed the master of %s", setNoMaster.Link)
		utils.ReplySuccess(w, r, msg, refreshedLink)
		return
	}

}
