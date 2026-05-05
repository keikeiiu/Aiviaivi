package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ailivili/internal/config"
	"ailivili/internal/db"
	"ailivili/internal/httpapi"
	"ailivili/internal/redis"
	"ailivili/internal/storage"
	"ailivili/internal/ws"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx := context.Background()

	sqlDB, err := db.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer sqlDB.Close()

	if err := db.ApplyMigrations(ctx, sqlDB, cfg.MigrationsDir); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	var redisClient *redis.Client
	if cfg.RedisURL != "" {
		redisClient, err = redis.New(ctx, cfg.RedisURL)
		if err != nil {
			log.Printf("redis: %v — continuing without cache", err)
			redisClient = nil
		} else {
			defer redisClient.Close()
			log.Printf("redis: connected to %s", cfg.RedisURL)
		}
	}

	hub := ws.NewHub()

	// Optional Redis pub/sub for cross-instance danmaku sync
	if redisClient != nil {
		rps := ws.NewRedisPubSub(redisClient.RDB(), hub)
		hub.SetRedisPubSub(rps)
		rps.Start(ctx)
	}

	var store storage.FileStore
	switch cfg.Storage {
	case "minio":
		useSSL := os.Getenv("MINIO_USE_SSL") == "true"
		s, err := storage.NewMinioStore(
			os.Getenv("MINIO_ENDPOINT"),
			os.Getenv("MINIO_ACCESS_KEY"),
			os.Getenv("MINIO_SECRET_KEY"),
			os.Getenv("MINIO_BUCKET"),
			cfg.StorageBaseURL,
			useSSL,
		)
		if err != nil {
			log.Fatalf("minio: %v", err)
		}
		store = s
	default:
		store = storage.NewLocalStore("uploads", cfg.StorageBaseURL)
	}

	deps := httpapi.Deps{
		DB:         sqlDB,
		JWTSecret:  cfg.JWTSecret,
		JWTExpires: cfg.JWTExpires,
		Hub:        hub,
		Store:      store,
		CORSOrigin: cfg.CORSOrigin,
	}
	if redisClient != nil {
		deps.Redis = redisClient.RDB()
	}

	// Periodically broadcast view counts to WebSocket clients
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			for _, videoID := range hub.ActiveRooms() {
				var count int64
				if err := sqlDB.QueryRowContext(ctx,
					`SELECT view_count FROM videos WHERE id = $1`, videoID,
				).Scan(&count); err == nil {
					hub.BroadcastViewCount(videoID, count)
				}
			}
		}
	}()

	srv := &http.Server{
		Addr:              cfg.Addr(),
		Handler:           httpapi.New(deps),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-stop:
	case err := <-errCh:
		log.Printf("server error: %v", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
