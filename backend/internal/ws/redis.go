package ws

import (
	"context"
	"encoding/json"
	"log"

	goredis "github.com/redis/go-redis/v9"
)

const danmakuChannelPrefix = "danmaku:"

// RedisPubSub bridges WebSocket danmaku broadcasts across server instances via Redis.
type RedisPubSub struct {
	rdb *goredis.Client
	hub *Hub
}

// NewRedisPubSub creates a Redis pub/sub adapter for the hub.
// Pass nil for rdb to disable cross-instance sync.
func NewRedisPubSub(rdb *goredis.Client, hub *Hub) *RedisPubSub {
	if rdb == nil {
		return nil
	}
	return &RedisPubSub{rdb: rdb, hub: hub}
}

// Start begins listening for danmaku messages from Redis.
// It spawns a goroutine that subscribes to all danmaku channels.
func (r *RedisPubSub) Start(ctx context.Context) {
	pubsub := r.rdb.PSubscribe(ctx, danmakuChannelPrefix+"*")
	ch := pubsub.Channel()

	go func() {
		defer pubsub.Close()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				r.handleRedisMessage(msg)
			}
		}
	}()

	log.Printf("ws: Redis pub/sub started for danmaku")
}

// Publish sends a danmaku message to Redis for other instances.
func (r *RedisPubSub) Publish(ctx context.Context, videoID string, msg Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	channel := danmakuChannelPrefix + videoID
	if err := r.rdb.Publish(ctx, channel, data).Err(); err != nil {
		log.Printf("ws redis publish error: %v", err)
	}
}

func (r *RedisPubSub) handleRedisMessage(msg *goredis.Message) {
	var dm Message
	if err := json.Unmarshal([]byte(msg.Payload), &dm); err != nil {
		log.Printf("ws redis unmarshal error: %v", err)
		return
	}

	// Broadcast to local clients (don't re-publish to Redis to avoid loops)
	r.hub.BroadcastToAll(dm.VideoID, dm)
}
