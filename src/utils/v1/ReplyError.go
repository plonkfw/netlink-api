package utilsv1

import (
	"encoding/json"
	"net/http"

	typesv1 "github.com/plonkfw/netlink-api/types/v1"
)

// ReplyError returns a formatted error response
func ReplyError(w http.ResponseWriter, r *http.Request, msg string, code string, err error) {
	// Setup some defaults
	status := "failed"
	message := ""

	if msg == "" {
		message = err.Error()
	} else {
		message = msg
	}

	errHeader := http.StatusBadRequest

	// Select proper error code based on message
	switch code {
	case "EEXISTS":
		errHeader = http.StatusConflict
		break

	case "ENOTPERMITTED":
		errHeader = http.StatusInternalServerError
		break

	case "ENOTIMPLEMENTED":
		errHeader = http.StatusNotImplemented
		break

	default:
		// If no error message is supplied, fall back to EEUNKNOWN
		code = "EEUNKNOWN"
		errHeader = http.StatusBadRequest
		break
	}

	// Prep the response
	response := typesv1.APIResponse{
		Status:  status,
		Code:    code,
		Message: message,
		Data:    err.Error(),
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		msg := "Error marshaling response"
		Log.Error().Err(err).Msg(msg)
		ReplyError(w, r, msg, "EPACKFAIL", err)
		return
	}

	// send the response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(errHeader)
	w.Write([]byte(jsonResponse))
}
