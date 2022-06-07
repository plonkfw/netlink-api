package link

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/plonkfw/netlink-api/utils"
	"github.com/vishvananda/netlink"
)

type linkByNameResponse struct {
	Status  string       `json:"status"`
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Data    netlink.Link `json:"data"`
}

// ByName retrieves a link by name
func ByName(w http.ResponseWriter, r *http.Request) {
	// Get the name of the link from query params - /v1/link/by-name?name=lo
	name := r.URL.Query().Get("name")

	// Lookup the link by name
	link, err := netlink.LinkByName(name)
	if err != nil {
		msg := fmt.Sprintf("Error querying link %s", name)
		utils.Log.Error().Err(err).Msg(msg)
		utils.ReplyError(w, r, msg, err)
		return
	}
	// Prep response
	msg := fmt.Sprintf("Found interface %s", name)
	response := linkByNameResponse{
		Status:  "success",
		Code:    "SUCCESS",
		Message: msg,
		Data:    link,
	}

	// JSON-ify the response
	jsonResponse, _ := json.MarshalIndent(response, "", "  ")

	// Send the response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonResponse))
}
