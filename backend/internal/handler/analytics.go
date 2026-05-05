package handler

import (
	"net/http"

	"ailivili/internal/middleware"
	"ailivili/internal/model"
	"ailivili/internal/response"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) AnalyticsOverview(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())

	overview, err := model.GetCreatorOverview(r.Context(), h.db, userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, 50080, "get analytics failed")
		return
	}

	response.OK(w, overview)
}

func (h *Handler) AnalyticsVideos(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())

	stats, err := model.GetCreatorVideoStats(r.Context(), h.db, userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, 50081, "get video stats failed")
		return
	}
	if stats == nil {
		stats = []model.VideoStats{}
	}

	response.OK(w, stats)
}

func (h *Handler) AnalyticsVideoDetail(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	videoID := chi.URLParam(r, "id")

	v, err := model.GetVideoByID(r.Context(), h.db, videoID)
	if err != nil {
		if err == model.ErrVideoNotFound {
			response.Error(w, http.StatusNotFound, 40402, "video not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, 50011, "get video failed")
		return
	}
	if v.UserID != userID {
		response.Error(w, http.StatusForbidden, 40301, "not your video")
		return
	}

	detail := model.VideoStatsDetail{
		VideoStats: model.VideoStats{
			ID:           v.ID,
			Title:        v.Title,
			CoverURL:     v.CoverURL,
			ViewCount:    v.ViewCount,
			LikeCount:    v.LikeCount,
			CommentCount: v.CommentCount,
			ShareCount:   v.ShareCount,
			PublishedAt:  v.PublishedAt,
			CreatedAt:    v.CreatedAt,
		},
	}

	response.OK(w, detail)
}
