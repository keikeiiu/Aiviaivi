package redis

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type Client struct {
	rdb *goredis.Client
}

func New(ctx context.Context, addr string) (*Client, error) {
	rdb := goredis.NewClient(&goredis.Options{
		Addr:        addr,
		DialTimeout: 5 * time.Second,
		PoolSize:    10,
		MinIdleConns: 2,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return &Client{rdb: rdb}, nil
}

func (c *Client) Close() error {
	return c.rdb.Close()
}

func (c *Client) RDB() *goredis.Client {
	return c.rdb
}

func (c *Client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}
