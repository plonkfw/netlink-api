package linkv1

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	utilsv1 "github.com/plonkfw/netlink-api/utils/v1"
	"github.com/vishvananda/netlink"
)

type setMaster struct {
	Name   string
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
		utilsv1.Log.Error().Err(err).Msg("Error reading body")
		utilsv1.ReplyError(w, r, "Error reading body", err)
		return
	}

	// Print the request to deubg stream
	utilsv1.Log.Debug().Msg(string(body))

	// Unpack the request
	if err := json.Unmarshal(body, &setMaster); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			utilsv1.Log.Error().Err(err).Msg("Error unmarshaling body")
			utilsv1.ReplyError(w, r, "Error unmarshaling body", err)
			return
		}
	}

	// If they supplied an interface name
	if setMaster.Master != "" && setMaster.Name != "" {
		// Look up the master link
		newMaster, err := netlink.LinkByName(setMaster.Master)
		if err != nil {
			msg := fmt.Sprintf("Error looking up master %s", setMaster.Master)
			utilsv1.Log.Error().Err(err).Msg(msg)
			utilsv1.ReplyError(w, r, msg, err)
			return
		}

		// Look up the child link
		newLink, err := netlink.LinkByName(setMaster.Name)
		if err != nil {
			msg := fmt.Sprintf("Error looking up link %s", setMaster.Name)
			utilsv1.Log.Error().Err(err).Msg(msg)
			utilsv1.ReplyError(w, r, msg, err)
			return
		}

		// Bind the child link to the master
		err = nil
		err = netlink.LinkSetMaster(newLink, newMaster)
		if err != nil {
			msg := fmt.Sprintf("Error binding link %s to master %s", setMaster.Name, setMaster.Master)
			utilsv1.Log.Error().Err(err).Msg(msg)
			utilsv1.ReplyError(w, r, msg, err)
			return
		}

		// Lookup the link by name
		refreshedLink, _ := netlink.LinkByName(setMaster.Name)
		refreshedMaster, _ := netlink.LinkByName(setMaster.Master)

		var responseData responseData

		responseData.Link = refreshedLink
		responseData.Master = refreshedMaster

		// Prep response
		msg := fmt.Sprintf("Successfully bound %s to master %s", setMaster.Name, setMaster.Master)
		utilsv1.ReplySuccess(w, r, msg, responseData)
		return
	}

	msg := fmt.Sprintf("Invalid paramaters %s %s", setMaster.Name, setMaster.Master)
	utilsv1.Log.Error().Err(err).Msg(msg)
	utilsv1.ReplyError(w, r, msg, err)
	return
}
