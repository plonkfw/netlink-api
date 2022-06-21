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

type setMTU struct {
	Link string
	MTU  int
}

// SetMTU sets the mtu of the given link device
func SetMTU(w http.ResponseWriter, r *http.Request) {
	var setMTU setMTU

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
	if err := json.Unmarshal(body, &setMTU); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			msg := "Error unmarshaling body"
			utilsv1.Log.Error().Err(err).Msg(msg)
			utilsv1.ReplyError(w, r, msg, "EUNPACKFAIL", err)
			return
		}
	}
	if setMTU.MTU != 0 && setMTU.Link != "" {
		link, err := netlink.LinkByName(setMTU.Link)
		if err != nil {
			msg := fmt.Sprintf("Error looking up link %s", setMTU.Link)
			utilsv1.Log.Error().Err(err).Msg(msg)
			utilsv1.ReplyError(w, r, msg, "ELOOKUPFAIL", err)
			return
		}

		err = nil
		err = netlink.LinkSetMTU(link, setMTU.MTU)
		if err != nil {
			msg := fmt.Sprintf("Error setting link %s mtu to %d", setMTU.Link, setMTU.MTU)
			utilsv1.Log.Error().Err(err).Msg(msg)
			utilsv1.ReplyError(w, r, msg, "EACTIONFAIL", err)
			return
		}

		refreshedLink, err := netlink.LinkByName(setMTU.Link)
		if err != nil {
			msg := fmt.Sprintf("Error refreshing info for link %s", setMTU.Link)
			utilsv1.Log.Error().Err(err).Msg(msg)
			utilsv1.ReplyError(w, r, msg, "ELOOKUPFAIL", err)
			return
		}

		// Prep response
		msg := fmt.Sprintf("Successfully set link %s MTU to %d", setMTU.Link, setMTU.MTU)
		utilsv1.ReplySuccess(w, r, msg, refreshedLink)
		return
	}
}
