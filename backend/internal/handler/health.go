package handler

import (
	"net/http"

	"ailivili/internal/response"
)

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	response.OK(w, map[string]any{"status": "ok"})
}
