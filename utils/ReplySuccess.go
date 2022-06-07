package utils

import (
	"encoding/json"
	"net/http"

	"github.com/plonkfw/netlink-api/types"
)

// ReplySuccess sends a formatted success reponse
func ReplySuccess(w http.ResponseWriter, r *http.Request, msg string, data interface{}) {
	// Prep response
	// dataString := string(data)
	response := types.APIResponse{
		Status:  "success",
		Code:    "SUCCESS",
		Message: msg,
		Data:    data,
	}

	// JSON-ify the response
	jsonResponse, _ := json.MarshalIndent(response, "", "  ")

	// Send the response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonResponse))
}
