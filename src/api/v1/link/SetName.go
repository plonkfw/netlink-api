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

type setName struct {
	Link string
	Name string
}

// SetName sets the name of the link device
func SetName(w http.ResponseWriter, r *http.Request) {
	var setName setName

	// Read the request
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		utilsv1.Log.Error().Err(err).Msg("Error reading body")
		utilsv1.ReplyError(w, r, "Error reading body", err)
		return
	}
	// Print the request to deubg stream
	utilsv1.Log.Debug().Msg(string(body))
	// Unpack the request
	if err := json.Unmarshal(body, &setName); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			utilsv1.Log.Error().Err(err).Msg("Error unmarshaling body")
			utilsv1.ReplyError(w, r, "Error unmarshaling body", err)
			return
		}
	}

	if setName.Name != "" && setName.Link != "" {
		link, err := netlink.LinkByName(setName.Link)
		if err != nil {
			msg := fmt.Sprintf("Error looking up link %s", setName.Link)
			utilsv1.Log.Error().Err(err).Msg(msg)
			utilsv1.ReplyError(w, r, msg, err)
			return
		}

		err = nil
		err = netlink.LinkSetName(link, setName.Name)
		if err != nil {
			msg := fmt.Sprintf("Error renaming link %s to %s", setName.Link, setName.Name)
			utilsv1.Log.Error().Err(err).Msg(msg)
			utilsv1.ReplyError(w, r, msg, err)
			return
		}

		refreshedLink, _ := netlink.LinkByName(setName.Name)

		// Prep response
		msg := fmt.Sprintf("Successfully renamed %s to %s", setName.Link, setName.Name)
		utilsv1.ReplySuccess(w, r, msg, refreshedLink)
		return
	}

	msg := fmt.Sprintf("Invalid paramaters %s %s", setName.Link, setName.Name)
	utilsv1.Log.Error().Err(err).Msg(msg)
	utilsv1.ReplyError(w, r, msg, err)
	return
}
