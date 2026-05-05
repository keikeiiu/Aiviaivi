package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"ailivili/internal/middleware"
	"ailivili/internal/model"
	"ailivili/internal/response"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) UsersMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, 40102, "missing bearer token")
		return
	}

	u, err := model.GetUserByID(r.Context(), h.db, userID)
	if err != nil {
		if err == model.ErrUserNotFound {
			response.Error(w, http.StatusNotFound, 40401, "user not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, 50005, "get user failed")
		return
	}

	response.OK(w, map[string]any{
		"id":         u.ID,
		"username":   u.Username,
		"email":      u.Email,
		"avatar_url": u.AvatarURL,
		"bio":        u.Bio,
		"role":       u.Role,
		"created_at": u.CreatedAt,
	})
}

func (h *Handler) UserProfile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	u, err := model.GetUserByID(r.Context(), h.db, id)
	if err != nil {
		if err == model.ErrUserNotFound {
			response.Error(w, http.StatusNotFound, 40401, "user not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, 50005, "get user failed")
		return
	}

	followerCount, _ := model.GetFollowerCount(r.Context(), h.db, id)
	followingCount, _ := model.GetFollowingCount(r.Context(), h.db, id)

	// Check if current user follows this user
	var isFollowing bool
	if userID, ok := middleware.UserIDFromContext(r.Context()); ok {
		isFollowing, _ = model.IsFollowing(r.Context(), h.db, userID, id)
	}

	response.OK(w, map[string]any{
		"id":             u.ID,
		"username":       u.Username,
		"avatar_url":     u.AvatarURL,
		"bio":            u.Bio,
		"role":           u.Role,
		"follower_count": followerCount,
		"following_count": followingCount,
		"is_following":   isFollowing,
		"created_at":     u.CreatedAt,
	})
}

func (h *Handler) UserUpdate(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	id := chi.URLParam(r, "id")

	if userID != id {
		response.Error(w, http.StatusForbidden, 40302, "cannot update other user's profile")
		return
	}

	var req struct {
		Bio       string `json:"bio"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, 40001, "invalid json")
		return
	}

	req.Bio = strings.TrimSpace(req.Bio)
	req.AvatarURL = strings.TrimSpace(req.AvatarURL)

	u, err := model.UpdateUser(r.Context(), h.db, id, req.Bio, req.AvatarURL)
	if err != nil {
		if err == model.ErrUserNotFound {
			response.Error(w, http.StatusNotFound, 40401, "user not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, 50006, "update user failed")
		return
	}

	response.OK(w, map[string]any{
		"id":         u.ID,
		"username":   u.Username,
		"avatar_url": u.AvatarURL,
		"bio":        u.Bio,
		"role":       u.Role,
	})
}

func (h *Handler) UserVideos(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	page, _ := parseIntQuery(r, "page")
	size, _ := parseIntQuery(r, "size")

	videos, total, err := model.ListVideos(r.Context(), h.db, model.VideoListParams{
		Page: page, Size: size, UserID: id,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, 50010, "list videos failed")
		return
	}

	items := make([]model.VideoDetail, len(videos))
	for i, v := range videos {
		items[i] = v.Detail(nil, nil)
	}

	response.WithPagination(w, items, pageOrDefault(page), sizeOrDefault(size), total)
}

func (h *Handler) UserFavorites(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	page, _ := parseIntQuery(r, "page")
	size, _ := parseIntQuery(r, "size")

	videos, total, err := model.ListFavorites(r.Context(), h.db, id, page, size)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, 50045, "list favorites failed")
		return
	}

	items := make([]model.VideoDetail, len(videos))
	for i, v := range videos {
		items[i] = v.Detail(nil, nil)
	}

	response.WithPagination(w, items, pageOrDefault(page), sizeOrDefault(size), total)
}
