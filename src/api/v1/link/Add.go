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
	Link string
	Type string
}

// Add creates a new network link - equivalent to `ip link add $i`
func Add(w http.ResponseWriter, r *http.Request) {
	// Prep our new link
	var linkAdd linkAdd

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
	if err := json.Unmarshal(body, &linkAdd); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			msg := "Error unmarshaling body"
			utilsv1.Log.Error().Err(err).Msg(msg)
			utilsv1.ReplyError(w, r, msg, "EUNPACKFAIL", err)
			return
		}
	}

	if linkAdd.Type != "" && linkAdd.Link != "" {
		// Call the appropriate function
		switch linkAdd.Type {
		// Basic bridge
		case "bridge":
			addBridge(w, r, linkAdd)
		// Dummy device
		case "dummy":
			addDummy(w, r, linkAdd)
		// Bail out
		default:
			msg := "Link type not implemented"
			err := errors.New(msg)
			utilsv1.ReplyError(w, r, "ENOTIMPLEMENTED", msg, err)
		}
		return
	}

	msg := fmt.Sprintf("Invalid paramaters %s %s", linkAdd.Link, linkAdd.Type)
	err = errors.New(msg)
	utilsv1.ReplyError(w, r, "ENOTIMPLEMENTED", msg, err)
}

// addBridge creates a basic linux bridge
func addBridge(w http.ResponseWriter, r *http.Request, linkAdd linkAdd) {
	// Setup link attributes
	linkAttrs := netlink.NewLinkAttrs()
	linkAttrs.Name = linkAdd.Link

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

	refreshedLink, err := netlink.LinkByName(bridge.Name)
	if err != nil {
		msg := fmt.Sprintf("Error looking up link %s", bridge.Name)
		utilsv1.Log.Error().Err(err).Msg(msg)
		utilsv1.ReplyError(w, r, msg, "ELOOKUPFAIL", err)
		return
	}
	msg := fmt.Sprintf("Successfully added bridge %s", bridge.Name)
	utilsv1.ReplySuccess(w, r, msg, refreshedLink)
	return
}

// addDummy creates a dummy device
func addDummy(w http.ResponseWriter, r *http.Request, linkAdd linkAdd) {
	// Setup link attributes
	linkAttrs := netlink.NewLinkAttrs()
	linkAttrs.Name = linkAdd.Link

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

	// Lookup the link by name
	refreshedLink, err := netlink.LinkByName(dummy.Name)
	if err != nil {
		msg := fmt.Sprintf("Error looking up link %s", dummy.Name)
		utilsv1.Log.Error().Err(err).Msg(msg)
		utilsv1.ReplyError(w, r, msg, "ELOOKUPFAIL", err)
		return
	}
	msg := fmt.Sprintf("Successfully added dummy %s", dummy.Name)
	utilsv1.ReplySuccess(w, r, msg, refreshedLink)
	return
}
