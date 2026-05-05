package ws

import (
	"context"
	"encoding/json"
	"log"
	"sync"
)

// Message represents a WebSocket message to/from clients.
type Message struct {
	Type      string `json:"type"`
	ID        int64  `json:"id,omitempty"`
	VideoID   string `json:"video_id,omitempty"`
	UserID    string `json:"user_id,omitempty"`
	Username  string `json:"username,omitempty"`
	Content   string `json:"content,omitempty"`
	VideoTime float64 `json:"video_time,omitempty"`
	Color     string `json:"color,omitempty"`
	FontSize  string `json:"font_size,omitempty"`
	Mode      string `json:"mode,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	Error     string `json:"error,omitempty"`
	Count     int64  `json:"count,omitempty"`
}

// A Room manages all clients watching the same video.
type Room struct {
	mu      sync.RWMutex
	clients map[*Client]struct{}
}

// Hub manages all video rooms.
type Hub struct {
	mu     sync.RWMutex
	rooms  map[string]*Room // keyed by videoID
	pubsub *RedisPubSub     // optional cross-instance sync
}

func NewHub() *Hub {
	return &Hub{rooms: make(map[string]*Room)}
}

// SetRedisPubSub attaches a Redis pub/sub adapter for cross-instance broadcast.
func (h *Hub) SetRedisPubSub(ps *RedisPubSub) {
	h.pubsub = ps
}

func (h *Hub) getOrCreateRoom(videoID string) *Room {
	h.mu.RLock()
	r, ok := h.rooms[videoID]
	h.mu.RUnlock()
	if ok {
		return r
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	// Double-check after acquiring write lock
	if r, ok = h.rooms[videoID]; ok {
		return r
	}
	r = &Room{clients: make(map[*Client]struct{})}
	h.rooms[videoID] = r
	return r
}

func (h *Hub) Join(videoID string, c *Client) {
	room := h.getOrCreateRoom(videoID)
	room.mu.Lock()
	room.clients[c] = struct{}{}
	room.mu.Unlock()
}

func (h *Hub) Leave(videoID string, c *Client) {
	h.mu.RLock()
	room, ok := h.rooms[videoID]
	h.mu.RUnlock()
	if !ok {
		return
	}
	room.mu.Lock()
	delete(room.clients, c)
	room.mu.Unlock()
}

// Broadcast sends a message to all clients in a video room except the sender.
func (h *Hub) Broadcast(videoID string, msg Message, sender *Client) {
	h.mu.RLock()
	room, ok := h.rooms[videoID]
	h.mu.RUnlock()
	if !ok {
		return
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("ws broadcast marshal: %v", err)
		return
	}

	room.mu.RLock()
	defer room.mu.RUnlock()
	for c := range room.clients {
		if c != sender {
			c.send <- data
		}
	}
}

// BroadcastToAll sends a message to all clients in a video room.
func (h *Hub) BroadcastToAll(videoID string, msg Message) {
	h.mu.RLock()
	room, ok := h.rooms[videoID]
	h.mu.RUnlock()
	if !ok {
		return
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("ws broadcast marshal: %v", err)
		return
	}

	room.mu.RLock()
	defer room.mu.RUnlock()
	for c := range room.clients {
		c.send <- data
	}
}

// PublishDanmaku broadcasts a danmaku to all local clients AND publishes to Redis
// for cross-instance sync (if Redis is configured).
func (h *Hub) PublishDanmaku(msg Message) {
	h.BroadcastToAll(msg.VideoID, msg)

	if h.pubsub != nil {
		h.pubsub.Publish(context.Background(), msg.VideoID, msg)
	}
}

// ActiveRooms returns all video IDs with connected clients.
func (h *Hub) ActiveRooms() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	ids := make([]string, 0, len(h.rooms))
	for id := range h.rooms {
		ids = append(ids, id)
	}
	return ids
}

// BroadcastViewCount sends a view count update to all clients in a room.
func (h *Hub) BroadcastViewCount(videoID string, count int64) {
	msg := Message{
		Type:    "view_count",
		VideoID: videoID,
		Count:   count,
	}
	h.BroadcastToAll(videoID, msg)
}

// RoomCount returns the number of clients in a video room.
func (h *Hub) RoomCount(videoID string) int {
	h.mu.RLock()
	room, ok := h.rooms[videoID]
	h.mu.RUnlock()
	if !ok {
		return 0
	}
	room.mu.RLock()
	defer room.mu.RUnlock()
	return len(room.clients)
}
