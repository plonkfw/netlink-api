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

type setMaster struct {
	Link   string
	Master string
}

type responseData struct {
	Link   netlink.Link `json:"link"`
	Master netlink.Link `json:"master"`
}

// SetMaster sets the parent bridge for an interface
func SetMaster(w http.ResponseWriter, r *http.Request) {
	// Prep our object
	var setMaster setMaster

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
	if err := json.Unmarshal(body, &setMaster); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			utils.Log.Error().Err(err).Msg("Error unmarshaling body")
			utils.ReplyError(w, r, "Error unmarshaling body", err)
			return
		}
	}

	// If they supplied an interface name
	if setMaster.Master != "" && setMaster.Link != "" {
		// Look up the master link
		newMaster, err := netlink.LinkByName(setMaster.Master)
		if err != nil {
			msg := fmt.Sprintf("Error looking up master %s", setMaster.Master)
			utils.Log.Error().Err(err).Msg(msg)
			utils.ReplyError(w, r, msg, err)
			return
		}

		// Look up the child link
		newLink, err := netlink.LinkByName(setMaster.Link)
		if err != nil {
			msg := fmt.Sprintf("Error looking up link %s", setMaster.Link)
			utils.Log.Error().Err(err).Msg(msg)
			utils.ReplyError(w, r, msg, err)
			return
		}

		// Bind the child link to the master
		err = nil
		err = netlink.LinkSetMaster(newLink, newMaster)
		if err != nil {
			msg := fmt.Sprintf("Error binding link %s to master %s", setMaster.Link, setMaster.Master)
			utils.Log.Error().Err(err).Msg(msg)
			utils.ReplyError(w, r, msg, err)
			return
		}

		// Lookup the link by name
		refreshedLink, _ := netlink.LinkByName(setMaster.Link)
		refreshedMaster, _ := netlink.LinkByName(setMaster.Master)

		var responseData responseData

		responseData.Link = refreshedLink
		responseData.Master = refreshedMaster

		// Prep response
		msg := fmt.Sprintf("Successfully bound %s to master %s", setMaster.Link, setMaster.Master)
		utils.ReplySuccess(w, r, msg, responseData)
		return
	}

	msg := fmt.Sprintf("Invalid paramaters %s %s", setMaster.Link, setMaster.Master)
	utils.Log.Error().Err(err).Msg(msg)
	utils.ReplyError(w, r, msg, err)
	return
}
