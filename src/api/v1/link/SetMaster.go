package linkv1

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	utilsv1 "github.com/plonkfw/netlink-api/utils/v1"
	"github.com/vishvananda/netlink"
)

type setMaster struct {
	Link   string
	Master string
}

type responseDataSetMaster struct {
	Link   netlink.Link `json:"link"`
	Master netlink.Link `json:"master"`
}

// SetMaster sets the parent bridge for an interface
func SetMaster(w http.ResponseWriter, r *http.Request) {
	// Prep our object
	var setMaster setMaster

	// Read in the request
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		msg := "Error reading body"
		utilsv1.Log.Error().Err(err).Msg(msg)
		utilsv1.ReplyError(w, r, msg, "EREADFAIL", err)
		return
	}

	// Print the request to deubg stream
	utilsv1.Log.Debug().Msg(string(body))

	// Unpack the request
	if err := json.Unmarshal(body, &setMaster); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			msg := "Error unmarshaling body"
			utilsv1.Log.Error().Err(err).Msg(msg)
			utilsv1.ReplyError(w, r, msg, "EUNPACKFAIL", err)
			return
		}
	}

	// If they supplied an interface name
	if setMaster.Master != "" && setMaster.Link != "" {
		// Look up the master link
		newMaster, err := netlink.LinkByName(setMaster.Master)
		if err != nil {
			msg := fmt.Sprintf("Error looking up master %s", setMaster.Master)
			utilsv1.Log.Error().Err(err).Msg(msg)
			utilsv1.ReplyError(w, r, msg, "ELOOKUPFAIL", err)
			return
		}

		// Look up the child link
		newLink, err := netlink.LinkByName(setMaster.Link)
		if err != nil {
			msg := fmt.Sprintf("Error looking up link %s", setMaster.Link)
			utilsv1.Log.Error().Err(err).Msg(msg)
			utilsv1.ReplyError(w, r, msg, "ELOOKUPFAIL", err)
			return
		}

		// Bind the child link to the master
		err = nil
		err = netlink.LinkSetMaster(newLink, newMaster)
		if err != nil {
			msg := fmt.Sprintf("Error binding link %s to master %s", setMaster.Link, setMaster.Master)
			utilsv1.Log.Error().Err(err).Msg(msg)
			utilsv1.ReplyError(w, r, msg, "EACTIONFAIL", err)
			return
		}

		// Lookup the link by name
		refreshedLink, err := netlink.LinkByName(setMaster.Link)
		if err != nil {
			msg := fmt.Sprintf("Error refreshing info for link %s", setMaster.Link)
			utilsv1.Log.Error().Err(err).Msg(msg)
			utilsv1.ReplyError(w, r, msg, "ELOOKUPFAIL", err)
			return
		}

		err = nil
		refreshedMaster, err := netlink.LinkByName(setMaster.Master)
		if err != nil {
			msg := fmt.Sprintf("Error refreshing info for link %s", setMaster.Master)
			utilsv1.Log.Error().Err(err).Msg(msg)
			utilsv1.ReplyError(w, r, msg, "ELOOKUPFAIL", err)
			return
		}

		var responseData responseDataSetMaster

		responseData.Link = refreshedLink
		responseData.Master = refreshedMaster

		// Prep response
		msg := fmt.Sprintf("Successfully bound %s to master %s", setMaster.Link, setMaster.Master)
		utilsv1.ReplySuccess(w, r, msg, responseData)
		return
	}

	// Invalid params
	msg := fmt.Sprintf("Invalid paramaters %s %s", setMaster.Link, setMaster.Master)
	err = errors.New(msg)
	utilsv1.Log.Error().Err(err).Msg(msg)
	utilsv1.ReplyError(w, r, msg, "EINVALIDPARAM", err)
	return
}
