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

type setDown struct {
	Name string
}

// SetDown disables link devices
func SetDown(w http.ResponseWriter, r *http.Request) {
	var setDown setDown

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
	if err := json.Unmarshal(body, &setDown); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			utils.Log.Error().Err(err).Msg("Error unmarshaling body")
			utils.ReplyError(w, r, "Error unmarshaling body", err)
			return
		}
	}

	if setDown.Name != "" {
		link, _ := netlink.LinkByName(setDown.Name)
		if err != nil {
			msg := fmt.Sprintf("Error looking up link %s", setDown.Name)
			utils.Log.Error().Err(err).Msg(msg)
			utils.ReplyError(w, r, msg, err)
			return
		}

		err = nil
		err = netlink.LinkSetDown(link)
		if err != nil {
			msg := fmt.Sprintf("Error downing link %s", setDown.Name)
			utils.Log.Error().Err(err).Msg(msg)
			utils.ReplyError(w, r, msg, err)
			return
		}

		// Lookup the link by name
		refreshedLink, _ := netlink.LinkByName(setDown.Name)

		// Prep response
		msg := fmt.Sprintf("Successfully downed interface %s", setDown.Name)
		utils.ReplySuccess(w, r, msg, refreshedLink)
		return
	}
}
