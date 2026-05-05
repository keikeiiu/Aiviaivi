package httpapi

import (
	"database/sql"
	"net/http"
	"time"

	"ailivili/internal/handler"
	"ailivili/internal/middleware"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

type Deps struct {
	DB         *sql.DB
	JWTSecret  string
	JWTExpires time.Duration
}

func New(deps Deps) http.Handler {
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(30 * time.Second))

	h := handler.New(handler.Deps{
		DB:         deps.DB,
		JWTSecret:  deps.JWTSecret,
		JWTExpires: deps.JWTExpires,
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", h.Health)
		r.Post("/auth/register", h.AuthRegister)
		r.Post("/auth/login", h.AuthLogin)
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireAuth(deps.JWTSecret))
			r.Get("/users/me", h.UsersMe)
		})
	})

	return r
}
