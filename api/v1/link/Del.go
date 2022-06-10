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

type linkDel struct {
	Name string
}

// Del removes link devices
func Del(w http.ResponseWriter, r *http.Request) {
	var linkDel linkDel

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
	if err := json.Unmarshal(body, &linkDel); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			utils.Log.Error().Err(err).Msg("Error unmarshaling body")
			utils.ReplyError(w, r, "Error unmarshaling body", err)
			return
		}
	}

	if linkDel.Name != "" {
		link, _ := netlink.LinkByName(linkDel.Name)
		if err != nil {
			msg := fmt.Sprintf("Error looking up link %s", linkDel.Name)
			utils.Log.Error().Err(err).Msg(msg)
			utils.ReplyError(w, r, msg, err)
			return
		}

		err = nil
		err = netlink.LinkDel(link)
		if err != nil {
			msg := fmt.Sprintf("Error removing link %s", linkDel.Name)
			utils.Log.Error().Err(err).Msg(msg)
			utils.ReplyError(w, r, msg, err)
			return
		}

		// Prep response
		msg := fmt.Sprintf("Successfully removed link %s", linkDel.Name)
		utils.ReplySuccess(w, r, msg, nil)
		return
	}
}
