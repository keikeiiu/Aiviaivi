package handler

import (
	"log"
	"net/http"
	"strings"

	"ailivili/internal/auth"
	"ailivili/internal/model"
	"ailivili/internal/ws"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) DanmakuWS(w http.ResponseWriter, r *http.Request) {
	videoID := chi.URLParam(r, "id")

	// Verify video exists
	if _, err := model.GetVideoByID(r.Context(), h.db, videoID); err != nil {
		http.Error(w, "video not found", http.StatusNotFound)
		return
	}

	// Optional auth via query param
	var userID, username string
	if tokenStr := r.URL.Query().Get("token"); tokenStr != "" {
		tok, claims, err := auth.ParseToken(strings.TrimSpace(tokenStr), h.jwtSecret)
		if err == nil && tok != nil && tok.Valid {
			if sub, ok := claims["sub"].(string); ok && sub != "" {
				u, err := model.GetUserByID(r.Context(), h.db, sub)
				if err == nil {
					userID = u.ID
					username = u.Username
				}
			}
		}
	}

	c, err := ws.NewClient(h.hub, w, r, videoID, userID, username)
	if err != nil {
		log.Printf("ws upgrade error: %v", err)
		return
	}

	_ = c // client is now managed by the hub (readPump/writePump goroutines)
}
