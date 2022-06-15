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

type setNoMaster struct {
	Name string
}

// SetNoMaster removes the master of the link device
func SetNoMaster(w http.ResponseWriter, r *http.Request) {
	var setNoMaster setNoMaster

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
	if err := json.Unmarshal(body, &setNoMaster); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			msg := "Error unmarshaling body"
			utilsv1.Log.Error().Err(err).Msg(msg)
			utilsv1.ReplyError(w, r, msg, "EUNPACKFAIL", err)
			return
		}
	}

	if setNoMaster.Name != "" {
		link, _ := netlink.LinkByName(setNoMaster.Name)
		if err != nil {
			msg := fmt.Sprintf("Error looking up link %s", setNoMaster.Name)
			utilsv1.Log.Error().Err(err).Msg(msg)
			utilsv1.ReplyError(w, r, msg, "ELOOKUPFAIL", err)
			return
		}

		err = nil
		err = netlink.LinkSetNoMaster(link)
		if err != nil {
			msg := fmt.Sprintf("Error removing master of link %s", setNoMaster.Name)
			utilsv1.Log.Error().Err(err).Msg(msg)
			utilsv1.ReplyError(w, r, msg, "EACTIONFAIL", err)
			return
		}

		// Lookup the link by name
		refreshedLink, _ := netlink.LinkByName(setNoMaster.Name)

		// Prep response
		msg := fmt.Sprintf("Successfully removed the master of %s", setNoMaster.Name)
		utilsv1.ReplySuccess(w, r, msg, refreshedLink)
		return
	}
	msg := fmt.Sprintf("Invalid paramaters %s", setNoMaster.Name)
	utilsv1.Log.Error().Err(err).Msg(msg)
	utilsv1.ReplyError(w, r, msg, "EINVALIDPARAM", err)
	return
}
