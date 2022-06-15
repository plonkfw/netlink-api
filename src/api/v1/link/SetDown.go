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

type setDown struct {
	Name string
}

// SetDown disables link devices
func SetDown(w http.ResponseWriter, r *http.Request) {
	var setDown setDown

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
	if err := json.Unmarshal(body, &setDown); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			msg := "Error unmarshaling body"
			utilsv1.Log.Error().Err(err).Msg(msg)
			utilsv1.ReplyError(w, r, msg, "EUNPACKFAIL", err)
			return
		}
	}

	if setDown.Name != "" {
		link, err := netlink.LinkByName(setDown.Name)
		if err != nil {
			msg := fmt.Sprintf("Error looking up link %s", setDown.Name)
			utilsv1.Log.Error().Err(err).Msg(msg)
			utilsv1.ReplyError(w, r, msg, "ELOOKUPFAIL", err)
			return
		}

		err = nil
		err = netlink.LinkSetDown(link)
		if err != nil {
			msg := fmt.Sprintf("Error downing link %s", setDown.Name)
			utilsv1.Log.Error().Err(err).Msg(msg)
			utilsv1.ReplyError(w, r, msg, "EACTIONFAIL", err)
			return
		}

		refreshedLink, err := netlink.LinkByName(setDown.Name)
		if err != nil {
			msg := fmt.Sprintf("Error refreshing info for link %s", setDown.Name)
			utilsv1.Log.Error().Err(err).Msg(msg)
			utilsv1.ReplyError(w, r, msg, "ELOOKUPFAIL", err)
			return
		}

		// Prep response
		msg := fmt.Sprintf("Successfully downed interface %s", setDown.Name)
		utilsv1.ReplySuccess(w, r, msg, refreshedLink)
		return
	}
	msg := fmt.Sprintf("Invalid paramaters %s", setDown.Name)
	utilsv1.Log.Error().Err(err).Msg(msg)
	utilsv1.ReplyError(w, r, msg, "EINVALIDPARAM", err)
}
