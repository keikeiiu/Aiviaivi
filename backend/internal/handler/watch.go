package handler

import (
	"encoding/json"
	"net/http"

	"ailivili/internal/middleware"
	"ailivili/internal/model"
	"ailivili/internal/response"

	"github.com/go-chi/chi/v5"
)

type watchRequest struct {
	Progress float64 `json:"progress"`
	Duration float64 `json:"duration"`
}

func (h *Handler) WatchRecord(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	videoID := chi.URLParam(r, "id")

	var req watchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, 40001, "invalid json")
		return
	}

	if req.Progress < 0 {
		response.Error(w, http.StatusBadRequest, 40005, "invalid progress")
		return
	}

	if err := model.RecordWatch(r.Context(), h.db, userID, videoID, req.Progress, req.Duration); err != nil {
		response.Error(w, http.StatusInternalServerError, 50060, "record watch failed")
		return
	}

	response.OK(w, nil)
}

func (h *Handler) WatchHistory(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	page, _ := parseIntQuery(r, "page")
	size, _ := parseIntQuery(r, "size")

	entries, total, err := model.GetWatchHistory(r.Context(), h.db, userID, page, size)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, 50061, "get watch history failed")
		return
	}

	response.WithPagination(w, entries, pageOrDefault(page), sizeOrDefault(size), total)
}
