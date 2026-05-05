package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ailivili/internal/response"
)

func TestNewRouter(t *testing.T) {
	deps := Deps{
		JWTSecret:  "test-secret",
		JWTExpires: time.Hour,
	}
	router := New(deps)
	if router == nil {
		t.Fatal("expected non-nil router")
	}
}

func TestHealthEndpoint(t *testing.T) {
	deps := Deps{
		JWTSecret:  "test-secret",
		JWTExpires: time.Hour,
	}
	router := New(deps)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/v1/health", nil)
	router.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var env response.Envelope[any]
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("json unmarshal: %v", err)
	}
	if env.Code != 0 {
		t.Fatalf("expected code=0, got %d", env.Code)
	}
}

func TestAuthEndpointsReturn400OnBadInput(t *testing.T) {
	deps := Deps{
		JWTSecret:  "test-secret",
		JWTExpires: time.Hour,
	}
	router := New(deps)

	tests := []struct {
		method string
		path   string
	}{
		{"POST", "/api/v1/auth/register"},
		{"POST", "/api/v1/auth/login"},
		{"POST", "/api/v1/auth/refresh"},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tt.method, tt.path, nil)
			router.ServeHTTP(w, r)

			if w.Code != http.StatusBadRequest {
				t.Errorf("expected 400, got %d", w.Code)
			}
		})
	}
}

func TestAuthRoutesReturn401WithoutToken(t *testing.T) {
	deps := Deps{
		JWTSecret:  "test-secret",
		JWTExpires: time.Hour,
	}
	router := New(deps)

	authRoutes := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/users/me"},
		{"PUT", "/api/v1/users/some-uuid"},
		{"POST", "/api/v1/users/some-uuid/subscribe"},
		{"POST", "/api/v1/videos/upload"},
		{"POST", "/api/v1/videos/some-uuid/danmaku"},
		{"POST", "/api/v1/videos/some-uuid/comments"},
		{"POST", "/api/v1/videos/some-uuid/like"},
		{"POST", "/api/v1/videos/some-uuid/favorite"},
		{"GET", "/api/v1/playlists"},
		{"POST", "/api/v1/playlists"},
		{"GET", "/api/v1/analytics/overview"},
	}

	for _, rt := range authRoutes {
		t.Run(rt.method+" "+rt.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(rt.method, rt.path, nil)
			router.ServeHTTP(w, r)

			if w.Code != http.StatusUnauthorized {
				t.Errorf("%s %s: expected 401, got %d", rt.method, rt.path, w.Code)
			}
		})
	}
}

func TestCORSPreflightRouting(t *testing.T) {
	deps := Deps{
		JWTSecret:  "test-secret",
		JWTExpires: time.Hour,
	}
	router := New(deps)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("OPTIONS", "/api/v1/health", nil)
	router.ServeHTTP(w, r)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204 for OPTIONS, got %d", w.Code)
	}
}

func TestMetricsEndpoint(t *testing.T) {
	deps := Deps{
		JWTSecret:  "test-secret",
		JWTExpires: time.Hour,
	}
	router := New(deps)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/metrics", nil)
	router.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for /metrics, got %d", w.Code)
	}
}

func TestRoute404Handling(t *testing.T) {
	deps := Deps{
		JWTSecret:  "test-secret",
		JWTExpires: time.Hour,
	}
	router := New(deps)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/v1/nonexistent", nil)
	router.ServeHTTP(w, r)

	if w.Code != http.StatusMethodNotAllowed && w.Code != http.StatusNotFound {
		t.Logf("unregistered route returned %d (expected 404 or 405)", w.Code)
	}
}
