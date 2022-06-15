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

// linkAdd for api/v1/link/Add.go
type linkAdd struct {
	Name string
	Type string
	MTU  int
}

// Add creates a new network link - equivalent to `ip link add $i`
func Add(w http.ResponseWriter, r *http.Request) {
	// Prep our new link
	var link linkAdd

	// Unpack the request
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
	if err := json.Unmarshal(body, &link); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			msg := "Error unmarshaling body"
			utilsv1.Log.Error().Err(err).Msg(msg)
			utilsv1.ReplyError(w, r, msg, "EUNPACKFAIL", err)
			return
		}
	}

	// Call the appropriate function
	switch link.Type {
	// Basic bridge
	case "bridge":
		addBridge(w, r, link)
	// Dummy device
	case "dummy":
		addDummy(w, r, link)
	// Bail out
	default:
		err := errors.New("Link type not implemented")
		utilsv1.ReplyError(w, r, "ENOTIMPLEMENTED", "Link type not implemented", err)
	}
}

// addBridge creates a basic linux bridge
func addBridge(w http.ResponseWriter, r *http.Request, link linkAdd) {
	// Setup link attributes
	linkAttrs := netlink.NewLinkAttrs()
	linkAttrs.Name = link.Name

	// Setup the netlink.Bridge struct
	bridge := &netlink.Bridge{
		LinkAttrs: linkAttrs,
	}

	// Create the bridge device
	err := netlink.LinkAdd(bridge)
	if err != nil {
		msg := fmt.Sprintf("Could not add bridge %s", bridge.Name)
		utilsv1.Log.Error().Err(err).Msg(msg)
		utilsv1.ReplyError(w, r, msg, "EACTIONFAIL", err)
		return
	}

	msg := fmt.Sprintf("Successfully added bridge %s", bridge.Name)
	utilsv1.ReplySuccess(w, r, msg, bridge)
}

// addDummy creates a dummy device
func addDummy(w http.ResponseWriter, r *http.Request, link linkAdd) {
	// Setup link attributes
	linkAttrs := netlink.NewLinkAttrs()
	linkAttrs.Name = link.Name

	// Setup the netlink.Dummy struct
	dummy := &netlink.Dummy{
		LinkAttrs: linkAttrs,
	}

	// Create the dummy device
	err := netlink.LinkAdd(dummy)
	if err != nil {
		msg := fmt.Sprintf("Could not add dummy %s", dummy.Name)
		utilsv1.Log.Error().Err(err).Msg(msg)
		utilsv1.ReplyError(w, r, msg, "EACTIONFAIL", err)
		return
	}

	msg := fmt.Sprintf("Successfully added dummy %s", dummy.Name)
	utilsv1.ReplySuccess(w, r, msg, dummy)
}
