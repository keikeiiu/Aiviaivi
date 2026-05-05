package handler

import (
	"net/http"

	"ailivili/internal/model"
	"ailivili/internal/response"
)

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	searchType := r.URL.Query().Get("type")
	page, _ := parseIntQuery(r, "page")
	size, _ := parseIntQuery(r, "size")

	if q == "" {
		response.Error(w, http.StatusBadRequest, 40002, "q is required")
		return
	}

	if searchType == "user" {
		users, total, err := model.SearchUsers(r.Context(), h.db, q, page, size)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, 50050, "search failed")
			return
		}
		items := make([]model.UserPublic, len(users))
		for i, u := range users {
			items[i] = u.Public()
		}
		response.WithPagination(w, items, pageOrDefault(page), sizeOrDefault(size), total)
		return
	}

	videos, total, err := model.SearchVideos(r.Context(), h.db, q, page, size)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, 50050, "search failed")
		return
	}
	items := make([]model.VideoDetail, len(videos))
	for i, v := range videos {
		items[i] = v.Detail(nil, nil)
	}
	response.WithPagination(w, items, pageOrDefault(page), sizeOrDefault(size), total)
}

func (h *Handler) TrendingFeed(w http.ResponseWriter, r *http.Request) {
	page, _ := parseIntQuery(r, "page")
	size, _ := parseIntQuery(r, "size")

	videos, total, err := model.ListVideos(r.Context(), h.db, model.VideoListParams{
		Page: page, Size: size, Sort: "trending",
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

func (h *Handler) CategoriesList(w http.ResponseWriter, r *http.Request) {
	cats, err := model.ListCategories(r.Context(), h.db)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, 50051, "list categories failed")
		return
	}
	if cats == nil {
		cats = []model.Category{}
	}
	response.OK(w, cats)
}
