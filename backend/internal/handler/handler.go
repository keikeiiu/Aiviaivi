package handler

import (
	"database/sql"
	"time"
)

type Deps struct {
	DB         *sql.DB
	JWTSecret  string
	JWTExpires time.Duration
}

type Handler struct {
	db         *sql.DB
	jwtSecret  string
	jwtExpires time.Duration
}

func New(deps Deps) *Handler {
	return &Handler{
		db:         deps.DB,
		jwtSecret:  deps.JWTSecret,
		jwtExpires: deps.JWTExpires,
	}
}
