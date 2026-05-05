package middleware

import (
	"database/sql"
	"net/http"

	"ailivili/internal/model"
	"ailivili/internal/response"
)

func RequireRole(db *sql.DB, roles ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := UserIDFromContext(r.Context())
			if !ok {
				response.Error(w, http.StatusUnauthorized, 40102, "missing bearer token")
				return
			}

			u, err := model.GetUserByID(r.Context(), db, userID)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, 40103, "invalid token")
				return
			}

			for _, role := range roles {
				if u.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			response.Error(w, http.StatusForbidden, 40303, "insufficient permissions")
		})
	}
}
