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

type setMasterByIndex struct {
	Name        string
	MasterIndex int
}

type responseDataSetMasterByIndex struct {
	Link   netlink.Link `json:"link"`
	Master netlink.Link `json:"master"`
}

// SetMasterByIndex sets the parent bridge for an interface
func SetMasterByIndex(w http.ResponseWriter, r *http.Request) {
	// Prep our object
	var setMasterByIndex setMasterByIndex

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
	if err := json.Unmarshal(body, &setMasterByIndex); err != nil {
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
	if setMasterByIndex.MasterIndex != 0 && setMasterByIndex.Name != "" {
		// Look up the master link
		newMaster, err := netlink.LinkByIndex(setMasterByIndex.MasterIndex)
		if err != nil {
			msg := fmt.Sprintf("Error looking up master %d", setMasterByIndex.MasterIndex)
			utilsv1.Log.Error().Err(err).Msg(msg)
			utilsv1.ReplyError(w, r, msg, "ELOOKUPFAIL", err)
			return
		}

		// Look up the child link
		newLink, err := netlink.LinkByName(setMasterByIndex.Name)
		if err != nil {
			msg := fmt.Sprintf("Error looking up link %s", setMasterByIndex.Name)
			utilsv1.Log.Error().Err(err).Msg(msg)
			utilsv1.ReplyError(w, r, msg, "ELOOKUPFAIL", err)
			return
		}

		// Bind the child link to the master
		err = nil
		err = netlink.LinkSetMaster(newLink, newMaster)
		if err != nil {
			msg := fmt.Sprintf("Error binding link %s to master %d", setMasterByIndex.Name, setMasterByIndex.MasterIndex)
			utilsv1.Log.Error().Err(err).Msg(msg)
			utilsv1.ReplyError(w, r, msg, "EACTIONFAIL", err)
			return
		}

		// Lookup the link by name
		refreshedLink, _ := netlink.LinkByName(setMasterByIndex.Name)
		refreshedMaster, _ := netlink.LinkByIndex(setMasterByIndex.MasterIndex)

		var responseData responseDataSetMasterByIndex

		responseData.Link = refreshedLink
		responseData.Master = refreshedMaster

		// Prep response
		msg := fmt.Sprintf("Successfully bound %s to master %d", setMasterByIndex.Name, setMasterByIndex.MasterIndex)
		utilsv1.ReplySuccess(w, r, msg, responseData)
		return
	}

	msg := fmt.Sprintf("Invalid paramaters %s %d", setMasterByIndex.Name, setMasterByIndex.MasterIndex)
	utilsv1.Log.Error().Err(err).Msg(msg)
	utilsv1.ReplyError(w, r, msg, "EINVALIDPARAM", err)
	return
}
