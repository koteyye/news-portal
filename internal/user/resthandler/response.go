package resthandler

import (
	"encoding/json"
	"net/http"
)

const (
	ctJSON = "application/json"
)

type responseErrMessage struct {
	Msg string `json:"msg"`
}

func (h *RESTHandler) mapErrToResponse(w http.ResponseWriter, statusCode int, err error) {
	payload, err := json.Marshal(&responseErrMessage{Msg: err.Error()})
	w.Header().Add("Content-Type", ctJSON)
	w.WriteHeader(statusCode)
	w.Write(payload)
}