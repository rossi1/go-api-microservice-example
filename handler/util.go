package handler

import (
	"encoding/json"
	"net/http"
)

type ErrResponse struct {
	Error string `json:"error"`
} //@name ErrorResponse

func ErrorResponse(w http.ResponseWriter, msg string, code int) {
	resp := ErrResponse{Error: msg}
	status := http.StatusInternalServerError

	switch code {
	case 404:
		status = http.StatusNotFound
	case 400:
		status = http.StatusBadRequest
	}

	Response(w, resp, status)
}

func Response(w http.ResponseWriter, res interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")

	content, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)

	w.Write(content)
}
