package addr

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/plonkfw/netlink-api/utils"
	"github.com/vishvananda/netlink"
)

type addAddressResponse struct {
	Status  string         `json:"status"`
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Data    []netlink.Addr `json:"data"`
}

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
		msg := fmt.Sprintf("Error reading body")
		utils.Log.Error().Err(err).Msg(msg)
		utils.ReplyError(w, r, msg, err)
		return
	}

	// Print the request to deubg stream
	utils.Log.Debug().Msg(string(body))

	// Unpack the request
	if err := json.Unmarshal(body, &addr); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			msg := fmt.Sprintf("Error unmarshaling body")
			utils.Log.Error().Err(err).Msg(msg)
			utils.ReplyError(w, r, msg, err)
			return
		}
	}

	// Lookup the link devices by name
	link, err := netlink.LinkByName(addr.Link)
	if err != nil {
		msg := fmt.Sprintf("Error looking up link %s", addr.Link)
		utils.Log.Error().Err(err).Msg(msg)
		utils.ReplyError(w, r, msg, err)
		return
	}

	// Parse the given address
	parsedAddress, err := netlink.ParseAddr(addr.Address)
	if err != nil {
		msg := fmt.Sprintf("Error parsing address %s", addr.Address)
		utils.Log.Error().Err(err).Msg(msg)
		utils.ReplyError(w, r, msg, err)
		return
	}

	// Reset err to nil...
	err = nil
	// Attempt to create the link device
	err = netlink.AddrAdd(link, parsedAddress)

	// If it fails send our error response
	if err != nil {
		msg := fmt.Sprintf("Error adding address %s to link %s", addr.Address, addr.Link)
		utils.Log.Error().Err(err).Msg(msg)
		utils.ReplyError(w, r, msg, err)
		return
	}

	addressList, _ := netlink.AddrList(link, 0)

	// Prep response
	msg := fmt.Sprintf("Successfully added address %s to link %s", addr.Address, addr.Link)
	utils.ReplySuccess(w, r, msg, addressList)
}
