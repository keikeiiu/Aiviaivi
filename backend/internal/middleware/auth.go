package middleware

import (
	"context"
	"net/http"
	"strings"

	"ailivili/internal/auth"
	"ailivili/internal/response"
)

type ctxKey int

const userIDKey ctxKey = 1

func UserIDFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(userIDKey)
	s, ok := v.(string)
	return s, ok && s != ""
}

func RequireAuth(jwtSecret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			tokenStr, ok := parseBearer(authHeader)
			if !ok {
				response.Error(w, http.StatusUnauthorized, 40102, "missing bearer token")
				return
			}

			tok, claims, err := auth.ParseToken(tokenStr, jwtSecret)
			if err != nil || tok == nil || !tok.Valid {
				response.Error(w, http.StatusUnauthorized, 40103, "invalid token")
				return
			}

			sub, ok := claims["sub"].(string)
			if !ok || strings.TrimSpace(sub) == "" {
				response.Error(w, http.StatusUnauthorized, 40103, "invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, sub)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func parseBearer(v string) (string, bool) {
	v = strings.TrimSpace(v)
	if v == "" {
		return "", false
	}
	parts := strings.SplitN(v, " ", 2)
	if len(parts) != 2 {
		return "", false
	}
	if !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}
	tok := strings.TrimSpace(parts[1])
	if tok == "" {
		return "", false
	}
	return tok, true
}
