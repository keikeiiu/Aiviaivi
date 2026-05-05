package handler

import (
	"context"
	"database/sql"
	"log"
	"time"

	"ailivili/internal/metrics"
	"ailivili/internal/model"
	"ailivili/internal/storage"
	"ailivili/internal/ws"
)

type Deps struct {
	DB         *sql.DB
	JWTSecret  string
	JWTExpires time.Duration
	Hub        *ws.Hub
	Store      storage.FileStore
}

type Handler struct {
	db         *sql.DB
	jwtSecret  string
	jwtExpires time.Duration
	hub        *ws.Hub
	store      storage.FileStore
}

func New(deps Deps) *Handler {
	h := &Handler{
		db:         deps.DB,
		jwtSecret:  deps.JWTSecret,
		jwtExpires: deps.JWTExpires,
		hub:        deps.Hub,
		store:      deps.Store,
	}

	// Set up the WebSocket danmaku message handler
	ws.OnMessage = func(c *ws.Client, msg ws.Message) {
		h.handleDanmakuMessage(c, msg)
	}

	return h
}

func (h *Handler) handleDanmakuMessage(c *ws.Client, msg ws.Message) {
	// Validate authentication
	if !c.IsAuthenticated() {
		c.Send(ws.MakeError("authentication required to send danmaku"))
		return
	}

	// Validate content
	if msg.Content == "" {
		c.Send(ws.MakeError("content is required"))
		return
	}

	// Apply defaults
	color := msg.Color
	if color == "" {
		color = "#FFFFFF"
	}
	fontSize := msg.FontSize
	if fontSize == "" {
		fontSize = "medium"
	}
	mode := msg.Mode
	if mode == "" {
		mode = "scroll"
	}

	// Validate enums
	if fontSize != "small" && fontSize != "medium" && fontSize != "large" {
		c.Send(ws.MakeError("invalid font_size"))
		return
	}
	if mode != "scroll" && mode != "top" && mode != "bottom" {
		c.Send(ws.MakeError("invalid mode"))
		return
	}

	// Persist danmaku (use background context since this runs in WebSocket goroutine)
	d, err := model.CreateDanmaku(context.Background(), h.db, c.VideoID(), c.UserID(), msg.Content, msg.VideoTime, color, fontSize, mode)
	if err != nil {
		log.Printf("ws danmaku persist error: %v", err)
		c.Send(ws.MakeError("failed to save danmaku"))
		return
	}

	metrics.IncDanmaku()

	// Build the broadcast message with persisted data
	// Use client's username (the CreateDanmaku query doesn't JOIN users)
	broadcast := ws.Message{
		Type:      "danmaku",
		ID:        d.ID,
		VideoID:   d.VideoID,
		UserID:    d.UserID,
		Username:  msg.Username,
		Content:   d.Content,
		VideoTime: d.VideoTime,
		Color:     d.Color,
		FontSize:  d.FontSize,
		Mode:      d.Mode,
		CreatedAt: d.CreatedAt.Format(time.RFC3339),
	}

	// Broadcast locally + publish to Redis for cross-instance sync
	h.hub.PublishDanmaku(broadcast)
}
