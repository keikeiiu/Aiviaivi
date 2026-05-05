package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"ailivili/internal/middleware"
	"ailivili/internal/model"
	"ailivili/internal/response"

	"github.com/go-chi/chi/v5"
)

type danmakuSendReq struct {
	Content   string  `json:"content"`
	VideoTime float64 `json:"video_time"`
	Color     string  `json:"color"`
	FontSize  string  `json:"font_size"`
	Mode      string  `json:"mode"`
}

func (h *Handler) DanmakuList(w http.ResponseWriter, r *http.Request) {
	videoID := chi.URLParam(r, "id")

	tStart, _ := strconv.ParseFloat(r.URL.Query().Get("t_start"), 64)
	tEnd, _ := strconv.ParseFloat(r.URL.Query().Get("t_end"), 64)

	list, err := model.GetDanmakuByTimeRange(r.Context(), h.db, videoID, tStart, tEnd)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, 50020, "get danmaku failed")
		return
	}
	if list == nil {
		list = []model.Danmaku{}
	}

	response.OK(w, list)
}

func (h *Handler) DanmakuSend(w http.ResponseWriter, r *http.Request) {
	videoID := chi.URLParam(r, "id")
	userID, _ := middleware.UserIDFromContext(r.Context())

	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
	var req danmakuSendReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, 40001, "invalid json")
		return
	}

	req.Content = strings.TrimSpace(req.Content)
	if req.Content == "" {
		response.Error(w, http.StatusBadRequest, 40002, "content is required")
		return
	}

	if req.Color == "" {
		req.Color = "#FFFFFF"
	}
	if req.FontSize == "" {
		req.FontSize = "medium"
	}
	if req.Mode == "" {
		req.Mode = "scroll"
	}

	// Validate enum values
	if req.FontSize != "small" && req.FontSize != "medium" && req.FontSize != "large" {
		response.Error(w, http.StatusBadRequest, 40005, "invalid font_size")
		return
	}
	if req.Mode != "scroll" && req.Mode != "top" && req.Mode != "bottom" {
		response.Error(w, http.StatusBadRequest, 40005, "invalid mode")
		return
	}

	d, err := model.CreateDanmaku(r.Context(), h.db, videoID, userID, req.Content, req.VideoTime, req.Color, req.FontSize, req.Mode)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, 50021, "create danmaku failed")
		return
	}

	response.OK(w, d)
}
