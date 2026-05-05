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

type playlistCreateReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
}

type playlistUpdateReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
}

type addVideoReq struct {
	VideoID string `json:"video_id"`
}

func (h *Handler) PlaylistList(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())

	playlists, err := model.ListUserPlaylists(r.Context(), h.db, userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, 50070, "list playlists failed")
		return
	}
	if playlists == nil {
		playlists = []model.Playlist{}
	}
	response.OK(w, playlists)
}

func (h *Handler) PlaylistCreate(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())

	var req playlistCreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, 40001, "invalid json")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		response.Error(w, http.StatusBadRequest, 40002, "name is required")
		return
	}

	p, err := model.CreatePlaylist(r.Context(), h.db, userID, req.Name, req.Description, req.IsPublic)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, 50071, "create playlist failed")
		return
	}

	response.OK(w, p)
}

func (h *Handler) PlaylistDetail(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	p, err := model.GetPlaylistByID(r.Context(), h.db, id)
	if err != nil {
		if err == model.ErrPlaylistNotFound {
			response.Error(w, http.StatusNotFound, 40407, "playlist not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, 50072, "get playlist failed")
		return
	}

	// If playlist is private, only the owner can view
	if !p.IsPublic {
		userID, ok := middleware.UserIDFromContext(r.Context())
		if !ok || userID != p.UserID {
			response.Error(w, http.StatusForbidden, 40304, "playlist is private")
			return
		}
	}

	videos, _ := model.GetPlaylistVideos(r.Context(), h.db, id)
	if videos == nil {
		videos = []model.PlaylistVideo{}
	}

	response.OK(w, map[string]any{
		"playlist": p,
		"videos":   videos,
	})
}

func (h *Handler) PlaylistUpdate(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	id := chi.URLParam(r, "id")

	existing, err := model.GetPlaylistByID(r.Context(), h.db, id)
	if err != nil {
		if err == model.ErrPlaylistNotFound {
			response.Error(w, http.StatusNotFound, 40407, "playlist not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, 50072, "get playlist failed")
		return
	}
	if existing.UserID != userID {
		response.Error(w, http.StatusForbidden, 40305, "not your playlist")
		return
	}

	var req playlistUpdateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, 40001, "invalid json")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		response.Error(w, http.StatusBadRequest, 40002, "name is required")
		return
	}

	p, err := model.UpdatePlaylist(r.Context(), h.db, id, req.Name, req.Description, req.IsPublic)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, 50073, "update playlist failed")
		return
	}

	response.OK(w, p)
}

func (h *Handler) PlaylistDelete(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	id := chi.URLParam(r, "id")

	existing, err := model.GetPlaylistByID(r.Context(), h.db, id)
	if err != nil {
		if err == model.ErrPlaylistNotFound {
			response.Error(w, http.StatusNotFound, 40407, "playlist not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, 50072, "get playlist failed")
		return
	}
	if existing.UserID != userID {
		response.Error(w, http.StatusForbidden, 40305, "not your playlist")
		return
	}

	if err := model.DeletePlaylist(r.Context(), h.db, id); err != nil {
		response.Error(w, http.StatusInternalServerError, 50074, "delete playlist failed")
		return
	}

	response.OK(w, nil)
}

func (h *Handler) PlaylistAddVideo(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	playlistID := chi.URLParam(r, "id")

	existing, err := model.GetPlaylistByID(r.Context(), h.db, playlistID)
	if err != nil {
		if err == model.ErrPlaylistNotFound {
			response.Error(w, http.StatusNotFound, 40407, "playlist not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, 50072, "get playlist failed")
		return
	}
	if existing.UserID != userID {
		response.Error(w, http.StatusForbidden, 40305, "not your playlist")
		return
	}

	var req addVideoReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, 40001, "invalid json")
		return
	}

	if req.VideoID == "" {
		response.Error(w, http.StatusBadRequest, 40002, "video_id is required")
		return
	}

	if err := model.AddVideoToPlaylist(r.Context(), h.db, playlistID, req.VideoID); err != nil {
		response.Error(w, http.StatusInternalServerError, 50075, "add video to playlist failed")
		return
	}

	response.OK(w, nil)
}

func (h *Handler) PlaylistRemoveVideo(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	playlistID := chi.URLParam(r, "id")
	videoID := chi.URLParam(r, "videoID")

	existing, err := model.GetPlaylistByID(r.Context(), h.db, playlistID)
	if err != nil {
		if err == model.ErrPlaylistNotFound {
			response.Error(w, http.StatusNotFound, 40407, "playlist not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, 50072, "get playlist failed")
		return
	}
	if existing.UserID != userID {
		response.Error(w, http.StatusForbidden, 40305, "not your playlist")
		return
	}

	if err := model.RemoveVideoFromPlaylist(r.Context(), h.db, playlistID, videoID); err != nil {
		if err == model.ErrPlaylistVideoNotFound {
			response.Error(w, http.StatusNotFound, 40408, "video not in playlist")
			return
		}
		response.Error(w, http.StatusInternalServerError, 50076, "remove video from playlist failed")
		return
	}

	response.OK(w, nil)
}
