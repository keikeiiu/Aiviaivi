package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"ailivili/internal/auth"
	"ailivili/internal/model"
	"ailivili/internal/response"
)

// maxBodySize limits request body to 1 MB.
const maxBodySize = 1 << 20

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) AuthRegister(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, 40001, "invalid json")
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Password = strings.TrimSpace(req.Password)

	if req.Username == "" || req.Email == "" || req.Password == "" {
		response.Error(w, http.StatusBadRequest, 40002, "missing fields")
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, 50001, "hash password failed")
		return
	}

	u, err := model.CreateUser(r.Context(), h.db, req.Username, req.Email, hash)
	if err != nil {
		if err == model.ErrUserExists {
			response.Error(w, http.StatusConflict, 40901, "user already exists")
			return
		}
		response.Error(w, http.StatusInternalServerError, 50002, "create user failed")
		return
	}

	token, err := auth.NewToken(u.ID, h.jwtSecret, h.jwtExpires)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, 50003, "create token failed")
		return
	}

	response.OK(w, map[string]any{
		"user": map[string]any{
			"id":       u.ID,
			"username": u.Username,
			"email":    u.Email,
		},
		"token": token,
	})
}

func (h *Handler) AuthLogin(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, 40001, "invalid json")
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Password = strings.TrimSpace(req.Password)
	if req.Email == "" || req.Password == "" {
		response.Error(w, http.StatusBadRequest, 40002, "missing fields")
		return
	}

	u, err := model.GetUserByEmail(r.Context(), h.db, req.Email)
	if err != nil {
		if err == model.ErrUserNotFound {
			response.Error(w, http.StatusUnauthorized, 40101, "invalid credentials")
			return
		}
		response.Error(w, http.StatusInternalServerError, 50004, "get user failed")
		return
	}

	if err := auth.CheckPassword(u.PasswordHash, req.Password); err != nil {
		response.Error(w, http.StatusUnauthorized, 40101, "invalid credentials")
		return
	}

	token, err := auth.NewToken(u.ID, h.jwtSecret, h.jwtExpires)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, 50003, "create token failed")
		return
	}

	response.OK(w, map[string]any{
		"user": map[string]any{
			"id":       u.ID,
			"username": u.Username,
			"email":    u.Email,
		},
		"token": token,
	})
}
