package utilsv1

import (
	"encoding/json"
	"net/http"

	typesv1 "github.com/plonkfw/netlink-api/types/v1"
)

// ReplySuccess sends a formatted success reponse
func ReplySuccess(w http.ResponseWriter, r *http.Request, msg string, data interface{}) {
	// Prep response
	response := typesv1.APIResponse{
		Status:  "success",
		Code:    "SUCCESS",
		Message: msg,
		Data:    data,
	}

	// JSON-ify the response
	jsonResponse, _ := json.Marshal(response)

	// Send the response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonResponse))
}
