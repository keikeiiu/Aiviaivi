package response

import (
	"encoding/json"
	"net/http"
)

type Envelope[T any] struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Data       T           `json:"data,omitempty"`
	Pagination interface{} `json:"pagination,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func OK(w http.ResponseWriter, data any) {
	WriteJSON(w, http.StatusOK, Envelope[any]{Code: 0, Message: "ok", Data: data})
}

func Error(w http.ResponseWriter, status int, code int, msg string) {
	WriteJSON(w, status, Envelope[any]{Code: code, Message: msg})
}
