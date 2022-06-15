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

type setPromiscOff struct {
	Link string
}

// SetPromiscOff disables promiscuous mode for the link
func SetPromiscOff(w http.ResponseWriter, r *http.Request) {
	var setPromiscOff setPromiscOff

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
	if err := json.Unmarshal(body, &setPromiscOff); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			msg := "Error unmarshaling body"
			utilsv1.Log.Error().Err(err).Msg(msg)
			utilsv1.ReplyError(w, r, msg, "EUNPACKFAIL", err)
			return
		}
	}

	if setPromiscOff.Link != "" {
		link, err := netlink.LinkByName(setPromiscOff.Link)
		if err != nil {
			msg := fmt.Sprintf("Error looking up link %s", setPromiscOff.Link)
			utilsv1.Log.Error().Err(err).Msg(msg)
			utilsv1.ReplyError(w, r, msg, "ELOOKUPFAIL", err)
			return
		}

		err = nil
		err = netlink.SetPromiscOff(link)
		if err != nil {
			msg := fmt.Sprintf("Error disabling promiscuous on link %s", setPromiscOff.Link)
			utilsv1.Log.Error().Err(err).Msg(msg)
			utilsv1.ReplyError(w, r, msg, "EACTIONFAIL", err)
			return
		}

		refreshedLink, _ := netlink.LinkByName(setPromiscOff.Link)

		// Prep response
		msg := fmt.Sprintf("Successfully disabled promiscuous on link %s", setPromiscOff.Link)
		utilsv1.ReplySuccess(w, r, msg, refreshedLink)
		return
	}

	msg := fmt.Sprintf("Invalid paramaters %s", setPromiscOff.Link)
	utilsv1.Log.Error().Err(err).Msg(msg)
	utilsv1.ReplyError(w, r, msg, "EINVALIDPARAMS", err)
	return
}
