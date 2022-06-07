package utils

import (
	"encoding/json"
	"net/http"

	"github.com/plonkfw/netlink-api/types"
)

// ReplyError returns a formatted error response
func ReplyError(w http.ResponseWriter, r *http.Request, msg string, err error) {
	// Setup some defaults
	status := "failed"
	code := "EEUNKNOWN"
	message := err.Error()
	if msg != "" {
		message = msg
	}

	// Select proper error code based on message
	switch err.Error() {
	case "file exists":
		code = "EEXISTS"
		// We set the http status header up here
		w.WriteHeader(http.StatusConflict)
		break

	case "operation not permitted":
		code = "ENOTPERMITTED"
		w.WriteHeader(http.StatusInternalServerError)
		break

	case "not implemented":
		code = "ENOTIMPLEMENTED"
		w.WriteHeader(http.StatusNotImplemented)
		break

	default:
		// If no error message is supplied, fall back to EEUNKNOWN
		break
	}

	// Prep the response
	// dataString := string(data)
	response := types.APIResponse{
		Status:  status,
		Code:    code,
		Message: message,
	}
	jsonResponse, _ := json.MarshalIndent(response, "", "  ")

	// send the response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte(jsonResponse))
}
