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

type commentCreateReq struct {
	Content  string  `json:"content"`
	ParentID *string `json:"parent_id"`
}

func (h *Handler) CommentList(w http.ResponseWriter, r *http.Request) {
	videoID := chi.URLParam(r, "id")
	page, _ := parseIntQuery(r, "page")
	size, _ := parseIntQuery(r, "size")
	sort := r.URL.Query().Get("sort")

	comments, total, err := model.ListComments(r.Context(), h.db, videoID, model.CommentListParams{
		Page: page, Size: size, Sort: sort,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, 50030, "list comments failed")
		return
	}

	response.WithPagination(w, comments, pageOrDefault(page), sizeOrDefault(size), total)
}

func (h *Handler) CommentCreate(w http.ResponseWriter, r *http.Request) {
	videoID := chi.URLParam(r, "id")
	userID, _ := middleware.UserIDFromContext(r.Context())

	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
	var req commentCreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, 40001, "invalid json")
		return
	}

	req.Content = strings.TrimSpace(req.Content)
	if req.Content == "" {
		response.Error(w, http.StatusBadRequest, 40002, "content is required")
		return
	}

	c, err := model.CreateComment(r.Context(), h.db, videoID, userID, req.Content, req.ParentID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, 50031, "create comment failed")
		return
	}

	// Update video comment count
	_, _ = h.db.ExecContext(r.Context(),
		`UPDATE videos SET comment_count = comment_count + 1 WHERE id = $1`, videoID)

	response.OK(w, c)
}

func (h *Handler) CommentDelete(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	commentID := chi.URLParam(r, "id")

	if err := model.DeleteComment(r.Context(), h.db, commentID, userID); err != nil {
		if err == model.ErrCommentNotFound {
			response.Error(w, http.StatusNotFound, 40403, "comment not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, 50032, "delete comment failed")
		return
	}

	response.OK(w, nil)
}

func parseIntQuery(r *http.Request, key string) (int, error) {
	return parseInt(r.URL.Query().Get(key))
}

func parseInt(s string) (int, error) {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, nil
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}
