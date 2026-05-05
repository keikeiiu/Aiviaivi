package handler

import (
	"net/http"

	"ailivili/internal/middleware"
	"ailivili/internal/model"
	"ailivili/internal/response"
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
		"id":       u.ID,
		"username": u.Username,
		"email":    u.Email,
	})
}
