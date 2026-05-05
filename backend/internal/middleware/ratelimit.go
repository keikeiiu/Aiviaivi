package middleware

import (
	"net/http"
	"strconv"
	"time"

	"ailivili/internal/response"

	goredis "github.com/redis/go-redis/v9"
)

func RateLimit(rdb *goredis.Client, limit int, window time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := "ratelimit:" + r.RemoteAddr + ":" + r.URL.Path
			ctx := r.Context()

			count, err := rdb.Incr(ctx, key).Result()
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			if count == 1 {
				rdb.Expire(ctx, key, window)
			}

			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(limit-int(count)))

			if count > int64(limit) {
				response.Error(w, http.StatusTooManyRequests, 42901, "rate limit exceeded")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
