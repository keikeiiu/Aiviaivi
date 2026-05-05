package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"ailivili/internal/metrics"
	"ailivili/internal/middleware"
	"ailivili/internal/model"
	"ailivili/internal/response"
	"ailivili/internal/transcoder"

	"github.com/go-chi/chi/v5"
)

func pageOrDefault(p int) int {
	if p < 1 {
		return 1
	}
	return p
}

func sizeOrDefault(s int) int {
	if s < 1 || s > 50 {
		return 20
	}
	return s
}

const uploadsDir = "uploads"

func (h *Handler) VideosList(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	sort := r.URL.Query().Get("sort")
	categoryStr := r.URL.Query().Get("category")

	var categoryID *int
	if cid, err := strconv.Atoi(categoryStr); err == nil {
		categoryID = &cid
	}

	videos, total, err := model.ListVideos(r.Context(), h.db, model.VideoListParams{
		Page: page, Size: size, Sort: sort, CategoryID: categoryID,
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

func (h *Handler) VideoDetail(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	v, err := model.GetVideoByID(r.Context(), h.db, id)
	if err != nil {
		if err == model.ErrVideoNotFound {
			response.Error(w, http.StatusNotFound, 40402, "video not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, 50011, "get video failed")
		return
	}

	qualities, _ := model.GetVideoQualities(r.Context(), h.db, id)

	var cat *model.Category
	if v.CategoryID != nil {
		cat, _ = model.GetCategoryByID(r.Context(), h.db, *v.CategoryID)
	}

	// Increment view count asynchronously (use background context since request ctx may cancel)
	go func() { _ = model.IncrementViewCount(context.Background(), h.db, id) }()

	response.OK(w, v.Detail(qualities, cat))
}

func (h *Handler) VideoUpload(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())

	r.Body = http.MaxBytesReader(w, r.Body, 100<<20) // 100 MB max
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		response.Error(w, http.StatusBadRequest, 40003, "invalid multipart form")
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	description := r.FormValue("description")
	tagsStr := r.FormValue("tags")
	categoryStr := r.FormValue("category_id")

	if title == "" {
		response.Error(w, http.StatusBadRequest, 40002, "title is required")
		return
	}

	var categoryID *int
	if cid, err := strconv.Atoi(categoryStr); err == nil {
		categoryID = &cid
	}

	var tags []string
	if tagsStr != "" {
		tags = strings.Split(tagsStr, ",")
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
		}
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.Error(w, http.StatusBadRequest, 40004, "file is required")
		return
	}
	defer file.Close()

	video, err := model.CreateVideo(r.Context(), h.db, userID, title, description, categoryID, tags)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, 50012, "create video failed")
		return
	}

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".mp4"
	}
	storagePath := filepath.Join("raw", video.ID+ext)

	_, err = h.store.Save(storagePath, file, header.Size)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, 50014, "save file failed")
		return
	}

	if err := model.UpdateVideoStatus(r.Context(), h.db, video.ID, "processing", 0, ""); err != nil {
		response.Error(w, http.StatusInternalServerError, 50015, "update status failed")
		return
	}

	hlsDir := filepath.Join(uploadsDir, "hls", video.ID)
	go h.transcodeVideo(video.ID, filepath.Join(h.store.Subdir(), storagePath), hlsDir)

	metrics.IncUpload()

	response.OK(w, map[string]any{
		"id":      video.ID,
		"title":   video.Title,
		"status":  "processing",
		"message": "video uploaded, processing will begin shortly",
	})
}

func (h *Handler) VideoUpdate(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	id := chi.URLParam(r, "id")

	v, err := model.GetVideoByID(r.Context(), h.db, id)
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

	var req struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, 40001, "invalid json")
		return
	}

	updated, err := model.UpdateVideoMeta(r.Context(), h.db, id, req.Title, req.Description, req.Tags)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, 50016, "update video failed")
		return
	}

	response.OK(w, updated.Detail(nil, nil))
}

func (h *Handler) VideoDelete(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())
	id := chi.URLParam(r, "id")

	v, err := model.GetVideoByID(r.Context(), h.db, id)
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

	if err := model.SoftDeleteVideo(r.Context(), h.db, id); err != nil {
		response.Error(w, http.StatusInternalServerError, 50017, "delete video failed")
		return
	}

	response.OK(w, nil)
}

func (h *Handler) VideoRelated(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	v, err := model.GetVideoByID(r.Context(), h.db, id)
	if err != nil {
		response.OK(w, []any{})
		return
	}
	if v.CategoryID == nil {
		response.OK(w, []any{})
		return
	}

	videos, total, err := model.ListVideos(r.Context(), h.db, model.VideoListParams{
		Page: 1, Size: 10, CategoryID: v.CategoryID,
	})
	if err != nil {
		response.OK(w, []any{})
		return
	}

	// Filter out current video
	items := make([]model.VideoDetail, 0)
	for _, rv := range videos {
		if rv.ID != id {
			items = append(items, rv.Detail(nil, nil))
		}
	}

	response.WithPagination(w, items, 1, 10, total-1)
}

func (h *Handler) transcodeVideo(videoID string, rawPath string, hlsDir string) {
	ctx := context.Background()
	result, err := transcoder.TranscodeToHLS(ctx, rawPath, hlsDir)
	if err != nil {
		log.Printf("transcode %s failed: %v", videoID, err)
		return
	}

	coverURL := ""
	if result.ThumbPath != "" {
		coverURL = "/" + result.ThumbPath
	}

	if err := model.UpdateVideoStatus(ctx, h.db, videoID, "published", result.Duration, coverURL); err != nil {
		log.Printf("update video %s status failed: %v", videoID, err)
		return
	}

	for _, q := range result.Qualities {
		manifestURL := "/" + filepath.ToSlash(q.ManifestURL)
		if _, err := model.CreateVideoQuality(ctx, h.db, videoID, q.Quality, manifestURL, q.Bitrate, q.FileSize); err != nil {
			log.Printf("create quality %s for video %s failed: %v", q.Quality, videoID, err)
		}
	}

	if err := model.PublishVideo(ctx, h.db, videoID); err != nil {
		log.Printf("publish video %s failed: %v", videoID, err)
		return
	}

	metrics.IncTranscoded()
	log.Printf("transcode %s complete: %d qualities, duration=%ds", videoID, len(result.Qualities), result.Duration)
}
