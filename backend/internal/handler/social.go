package handler

import (
	"net/http"

	"ailivili/internal/middleware"
	"ailivili/internal/model"
	"ailivili/internal/response"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) VideoLike(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	videoID := chi.URLParam(r, "id")

	if err := model.LikeVideo(r.Context(), h.db, userID, videoID); err != nil {
		response.Error(w, http.StatusInternalServerError, 50040, "like failed")
		return
	}
	response.OK(w, nil)
}

func (h *Handler) VideoUnlike(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	videoID := chi.URLParam(r, "id")

	if err := model.UnlikeVideo(r.Context(), h.db, userID, videoID); err != nil {
		if err == model.ErrNotExists {
			response.Error(w, http.StatusNotFound, 40404, "like not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, 50041, "unlike failed")
		return
	}
	response.OK(w, nil)
}

func (h *Handler) VideoFavorite(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	videoID := chi.URLParam(r, "id")

	if err := model.FavoriteVideo(r.Context(), h.db, userID, videoID); err != nil {
		response.Error(w, http.StatusInternalServerError, 50042, "favorite failed")
		return
	}
	response.OK(w, nil)
}

func (h *Handler) VideoUnfavorite(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	videoID := chi.URLParam(r, "id")

	if err := model.UnfavoriteVideo(r.Context(), h.db, userID, videoID); err != nil {
		if err == model.ErrNotExists {
			response.Error(w, http.StatusNotFound, 40405, "favorite not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, 50043, "unfavorite failed")
		return
	}
	response.OK(w, nil)
}

func (h *Handler) UserSubscribe(w http.ResponseWriter, r *http.Request) {
	followerID, _ := middleware.UserIDFromContext(r.Context())
	creatorID := chi.URLParam(r, "id")

	if err := model.FollowUser(r.Context(), h.db, followerID, creatorID); err != nil {
		response.Error(w, http.StatusBadRequest, 40006, err.Error())
		return
	}
	response.OK(w, nil)
}

func (h *Handler) UserUnsubscribe(w http.ResponseWriter, r *http.Request) {
	followerID, _ := middleware.UserIDFromContext(r.Context())
	creatorID := chi.URLParam(r, "id")

	if err := model.UnfollowUser(r.Context(), h.db, followerID, creatorID); err != nil {
		if err == model.ErrNotExists {
			response.Error(w, http.StatusNotFound, 40406, "follow not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, 50044, "unfollow failed")
		return
	}
	response.OK(w, nil)
}
