package addrv1

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	utilsv1 "github.com/plonkfw/netlink-api/utils/v1"
	"github.com/vishvananda/netlink"
)

// AddrAdd for api/v1/addr/Add.go
type addrAdd struct {
	Link    string
	Address string
}

// Add adds an address to a link device
func Add(w http.ResponseWriter, r *http.Request) {
	// Prep our new address
	var addr addrAdd

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
	if err := json.Unmarshal(body, &addr); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			msg := "Error unmarshaling body"
			utilsv1.Log.Error().Err(err).Msg(msg)
			utilsv1.ReplyError(w, r, msg, "EUNPACKFAIL", err)
			return
		}
	}

	// Lookup the link devices by name
	link, err := netlink.LinkByName(addr.Link)
	if err != nil {
		msg := fmt.Sprintf("Error looking up link %s", addr.Link)
		utilsv1.Log.Error().Err(err).Msg(msg)
		utilsv1.ReplyError(w, r, msg, "ELOOKUPFAIL", err)
		return
	}

	// Parse the given address
	parsedAddress, err := netlink.ParseAddr(addr.Address)
	if err != nil {
		msg := fmt.Sprintf("Error parsing address %s", addr.Address)
		utilsv1.Log.Error().Err(err).Msg(msg)
		utilsv1.ReplyError(w, r, msg, "EPARSEFAIL", err)
		return
	}

	// Reset err to nil...
	err = nil
	// Attempt to create the link device
	err = netlink.AddrAdd(link, parsedAddress)

	// If it fails send our error response
	if err != nil {
		msg := fmt.Sprintf("Error adding address %s to link %s", addr.Address, addr.Link)
		utilsv1.Log.Error().Err(err).Msg(msg)
		utilsv1.ReplyError(w, r, msg, "EACTIONFAIL", err)
		return
	}

	// Get address info
	addressList, err := netlink.AddrList(link, 0)
	if err != nil {
		msg := fmt.Sprintf("Error refreshing info for link %s", addr.Link)
		utilsv1.Log.Error().Err(err).Msg(msg)
		utilsv1.ReplyError(w, r, msg, "ELOOKUPFAIL", err)
		return
	}

	// Prep response
	msg := fmt.Sprintf("Successfully added address %s to link %s", addr.Address, addr.Link)
	utilsv1.ReplySuccess(w, r, msg, addressList)
}
