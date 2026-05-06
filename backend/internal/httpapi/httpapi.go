package httpapi

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"ailivili/internal/handler"
	"ailivili/internal/middleware"
	"ailivili/internal/storage"
	"ailivili/internal/ws"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	goredis "github.com/redis/go-redis/v9"
)

type Deps struct {
	DB         *sql.DB
	JWTSecret  string
	JWTExpires time.Duration
	Redis      *goredis.Client  // optional
	Hub        *ws.Hub
	Store      storage.FileStore
	CORSOrigin string
}

func New(deps Deps) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.PrometheusMetrics)
	r.Use(middleware.CORS(deps.CORSOrigin))
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	// Optional rate limiting via Redis
	if deps.Redis != nil {
		r.Use(middleware.RateLimit(deps.Redis, 100, time.Minute))
	}

	h := handler.New(handler.Deps{
		DB:         deps.DB,
		JWTSecret:  deps.JWTSecret,
		JWTExpires: deps.JWTExpires,
		Hub:        deps.Hub,
		Store:      deps.Store,
	})

	r.Route("/api/v1", func(r chi.Router) {
		// Apply timeout to REST routes (not WebSocket)
		r.Use(chimw.Timeout(30 * time.Second))
		// Public
		r.Get("/health", h.Health)
		r.Post("/auth/register", h.AuthRegister)
		r.Post("/auth/login", h.AuthLogin)
		r.Post("/auth/refresh", h.AuthRefresh)

		// Categories
		r.Get("/categories", h.CategoriesList)

		// Search
		r.Get("/search", h.Search)

		// Videos — public
		r.Get("/videos", h.VideosList)
		r.Get("/videos/{id}", h.VideoDetail)
		r.Get("/videos/{id}/related", h.VideoRelated)
		r.Get("/videos/{id}/danmaku", h.DanmakuList)
		r.Get("/videos/{id}/comments", h.CommentList)

		// Feed
		r.Get("/feed/trending", h.TrendingFeed)

		// Users — public profile
		r.Get("/users/{id}", h.UserProfile)
		r.Get("/users/{id}/videos", h.UserVideos)
		r.Get("/users/{id}/favorites", h.UserFavorites)

		// Playlists — public detail (private playlists rejected in handler)
		r.Get("/playlists/{id}", h.PlaylistDetail)

		// Authenticated routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireAuth(deps.JWTSecret))

			// User self
			r.Get("/users/me", h.UsersMe)
			r.Put("/users/{id}", h.UserUpdate)
			r.Post("/users/{id}/subscribe", h.UserSubscribe)
			r.Delete("/users/{id}/subscribe", h.UserUnsubscribe)

			// Videos — authenticated actions
			r.Post("/videos/upload", h.VideoUpload)
			r.Put("/videos/{id}", h.VideoUpdate)
			r.Delete("/videos/{id}", h.VideoDelete)

			// Danmaku
			r.Post("/videos/{id}/danmaku", h.DanmakuSend)

			// Comments
			r.Post("/videos/{id}/comments", h.CommentCreate)
			r.Delete("/comments/{id}", h.CommentDelete)

			// Social
			r.Post("/videos/{id}/like", h.VideoLike)
			r.Delete("/videos/{id}/like", h.VideoUnlike)
			r.Post("/videos/{id}/favorite", h.VideoFavorite)
			r.Delete("/videos/{id}/favorite", h.VideoUnfavorite)

			// Watch history
			r.Post("/videos/{id}/watch", h.WatchRecord)
			r.Get("/users/me/history", h.WatchHistory)

			// Playlists
			r.Get("/playlists", h.PlaylistList)
			r.Post("/playlists", h.PlaylistCreate)
			r.Put("/playlists/{id}", h.PlaylistUpdate)
			r.Delete("/playlists/{id}", h.PlaylistDelete)
			r.Post("/playlists/{id}/videos", h.PlaylistAddVideo)
			r.Delete("/playlists/{id}/videos/{videoID}", h.PlaylistRemoveVideo)

			// Analytics (creator)
			r.Get("/analytics/overview", h.AnalyticsOverview)
			r.Get("/analytics/videos", h.AnalyticsVideos)
			r.Get("/analytics/videos/{id}", h.AnalyticsVideoDetail)
		})
	})

	// WebSocket routes (outside timeout middleware to keep connections alive)
	r.Get("/api/v1/videos/{id}/danmaku/ws", h.DanmakuWS)

	// Prometheus metrics endpoint
	r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})

	// Serve HLS and uploaded files
	fs := http.FileServer(http.Dir("uploads"))
	r.Get("/uploads/*", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/uploads")
		fs.ServeHTTP(w, r)
	})

	return r
}
