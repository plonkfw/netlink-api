package link

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/plonkfw/netlink-api/types"
	"github.com/plonkfw/netlink-api/utils"
	"github.com/vishvananda/netlink"
)

type addBridgeResponse struct {
	Status  string          `json:"status"`
	Code    string          `json:"code"`
	Message string          `json:"message"`
	Data    *netlink.Bridge `json:"data"`
}

type addDummyResponse struct {
	Status  string         `json:"status"`
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Data    *netlink.Dummy `json:"data"`
}

// Add creates a new network link - equivalent to `ip link add $i`
func Add(w http.ResponseWriter, r *http.Request) {
	// Prep our new link
	var link types.LinkAdd

	// Unpack the request
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		utils.Log.Error().Err(err).Msg("Error reading body")
		utils.ReplyError(w, r, "Error reading body", err)
		return
	}

	// Print the request to deubg stream
	utils.Log.Debug().Msg(string(body))

	// Unpack the request
	if err := json.Unmarshal(body, &link); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			utils.Log.Error().Err(err).Msg("Error unmarshaling body")
			utils.ReplyError(w, r, "Error unmarshaling body", err)
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
		err := errors.New("not implemented")
		utils.ReplyError(w, r, "not implemented", err)
	}
}

// addBridge creates a basic linux bridge
func addBridge(w http.ResponseWriter, r *http.Request, link types.LinkAdd) {
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
		utils.Log.Error().Err(err).Msg(msg)
		utils.ReplyError(w, r, msg, err)
		return
	}

	// Prep response
	response := addBridgeResponse{
		Status:  "success",
		Code:    "SUCCESS",
		Message: fmt.Sprintf("Successfully added bridge %s", bridge.Name),
		Data:    bridge,
	}

	// JSON-ify the response
	jsonResponse, _ := json.MarshalIndent(response, "", "  ")

	// Send the response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonResponse))
}

// addDummy creates a dummy device
func addDummy(w http.ResponseWriter, r *http.Request, link types.LinkAdd) {
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
		utils.Log.Error().Err(err).Msg(msg)
		utils.ReplyError(w, r, msg, err)
		return
	}

	// Prep response
	response := addDummyResponse{
		Status:  "success",
		Code:    "SUCCESS",
		Message: fmt.Sprintf("Successfully added dummy %s", dummy.Name),
		Data:    dummy,
	}

	// JSON-ify the response
	jsonResponse, _ := json.MarshalIndent(response, "", "  ")

	// Send the response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonResponse))
}
