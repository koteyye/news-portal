package restresponser

import (
	"encoding/json"
	"net/http"
)

const (
	CtJSON = "application/json"
)

// ResponseOptions опции ответа
type ResponseOptions struct {
	StatusCode int
	Err error
	ContentType string
}

type responseErrMessage struct {
	Msg string `json:"msg"`
}

// MapErrToResponse маппит ошибку в HTTP ответ
func MapErrToResponse(w http.ResponseWriter, options *ResponseOptions) {
	w.Header().Set("Content-Type", options.ContentType)
	w.WriteHeader(options.StatusCode)
	if options.Err != nil {
		json.NewEncoder(w).Encode(responseErrMessage{Msg: options.Err.Error()})
	}
}
