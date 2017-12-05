package health

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type jsonStatus struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// NewBasicHandlerFunc will return an `http.HandlerFunc` that will write `ok`
// string + `http.StatusOK` to `rw`` if `h.Failed()` returns `false`;
// returns `error` + `http.StatusInternalServerError` if `h.Failed()` returns `true`.
func NewBasicHandlerFunc(h IHealth) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		status := http.StatusOK
		body := "ok"

		if h.Failed() {
			status = http.StatusInternalServerError
			body = "failed"
		}

		rw.WriteHeader(status)
		rw.Write([]byte(body))
	})
}

// NewJSONHandlerFunc will return an `http.HandlerFunc` that will marshal and
// write the contents of `h.StateMapInterface()` to `rw` and set status code to
//  `http.StatusOK` if `h.Failed()` is `false` OR set status code to
// `http.StatusInternalServerError` if `h.Failed` is `true`.
func NewJSONHandlerFunc(h IHealth) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		states, failed, err := h.State()
		if err != nil {
			writeJSONStatus(rw, "error", fmt.Sprintf("Unable to fetch states: %v", err), http.StatusOK)
			return
		}

		msg := "ok"
		statusCode := http.StatusOK

		if failed {
			msg = "failed"
			statusCode = http.StatusInternalServerError
		}

		fullBody := map[string]interface{}{
			"status":  msg,
			"details": states,
		}

		data, err := json.Marshal(fullBody)
		if err != nil {
			writeJSONStatus(rw, "error", fmt.Sprintf("Failed to marshal state data: %v", err), http.StatusOK)
			return
		}

		writeJSONResponse(rw, statusCode, data)
	})
}

func writeJSONStatus(rw http.ResponseWriter, status, message string, statusCode int) {
	jsonData, _ := json.Marshal(&jsonStatus{
		Message: message,
		Status:  status,
	})

	writeJSONResponse(rw, statusCode, jsonData)
}

func writeJSONResponse(rw http.ResponseWriter, statusCode int, content []byte) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	rw.Write(content)
}
