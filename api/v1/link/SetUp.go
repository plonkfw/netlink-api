package linkv1

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/plonkfw/netlink-api/utils"
	"github.com/vishvananda/netlink"
)

type setUp struct {
	Name string
}

// SetUp enables link devices
func SetUp(w http.ResponseWriter, r *http.Request) {
	var setUp setUp

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
	if err := json.Unmarshal(body, &setUp); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			utils.Log.Error().Err(err).Msg("Error unmarshaling body")
			utils.ReplyError(w, r, "Error unmarshaling body", err)
			return
		}
	}

	if setUp.Name != "" {
		link, _ := netlink.LinkByName(setUp.Name)
		if err != nil {
			msg := fmt.Sprintf("Error looking up link %s", setUp.Name)
			utils.Log.Error().Err(err).Msg(msg)
			utils.ReplyError(w, r, msg, err)
			return
		}

		err = nil
		err = netlink.LinkSetUp(link)
		if err != nil {
			msg := fmt.Sprintf("Error bringing up link %s", setUp.Name)
			utils.Log.Error().Err(err).Msg(msg)
			utils.ReplyError(w, r, msg, err)
			return
		}

		// Lookup the link by name
		refreshedLink, _ := netlink.LinkByName(setUp.Name)

		// Prep response
		msg := fmt.Sprintf("Successfully brought up interface %s", setUp.Name)
		utils.ReplySuccess(w, r, msg, refreshedLink)
		return
	}
}
