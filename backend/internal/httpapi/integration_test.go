//go:build integration
// +build integration

package httpapi

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"ailivili/internal/db"
	"ailivili/internal/response"
	"ailivili/internal/ws"
)

// setupIntegrationDB connects to the test database and runs migrations.
func setupIntegrationDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL not set — skipping integration test")
	}

	ctx := context.Background()
	sqlDB, err := db.Open(ctx, dsn)
	if err != nil {
		t.Fatalf("db open: %v", err)
	}

	if err := db.ApplyMigrations(ctx, sqlDB, "../../migrations"); err != nil {
		sqlDB.Close()
		t.Fatalf("migrations: %v", err)
	}

	return sqlDB
}

func newTestRouter(t *testing.T, sqlDB *sql.DB) http.Handler {
	t.Helper()
	return New(Deps{
		DB:         sqlDB,
		JWTSecret:  "integration-test-secret",
		JWTExpires: time.Hour,
		Hub:        ws.NewHub(),
	})
}

func registerAndLogin(t *testing.T, router http.Handler, username string) string {
	t.Helper()

	// Register
	body := fmt.Sprintf(`{"username":"%s","email":"%s@test.com","password":"pass"}`, username, username)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, r)

	var env response.Envelope[any]
	json.Unmarshal(w.Body.Bytes(), &env)
	if env.Code != 0 {
		t.Fatalf("register %s: code=%d msg=%s", username, env.Code, env.Message)
	}
	data := env.Data.(map[string]any)
	return data["token"].(string)
}

func login(t *testing.T, router http.Handler, username string) string {
	t.Helper()
	body := fmt.Sprintf(`{"email":"%s@test.com","password":"pass"}`, username)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, r)

	var env response.Envelope[any]
	json.Unmarshal(w.Body.Bytes(), &env)
	if env.Code != 0 {
		t.Fatalf("login %s: code=%d msg=%s", username, env.Code, env.Message)
	}
	data := env.Data.(map[string]any)
	return data["token"].(string)
}

func authRequest(method, path, token string, body io.Reader) *http.Request {
	r := httptest.NewRequest(method, path, body)
	if body != nil {
		r.Header.Set("Content-Type", "application/json")
	}
	r.Header.Set("Authorization", "Bearer "+token)
	return r
}

// --- Tests ---

func TestIntegrationHealth(t *testing.T) {
	sqlDB := setupIntegrationDB(t)
	defer sqlDB.Close()
	router := newTestRouter(t, sqlDB)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/v1/health", nil)
	router.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestIntegrationRegisterAndLogin(t *testing.T) {
	sqlDB := setupIntegrationDB(t)
	defer sqlDB.Close()
	router := newTestRouter(t, sqlDB)

	username := fmt.Sprintf("inttest_%d", time.Now().UnixNano())

	// Register
	body := fmt.Sprintf(`{"username":"%s","email":"%s@test.com","password":"testpass"}`, username, username)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatalf("register: expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var env response.Envelope[any]
	json.Unmarshal(w.Body.Bytes(), &env)
	data := env.Data.(map[string]any)
	if data["token"] == nil || data["refresh_token"] == nil {
		t.Fatal("expected token and refresh_token in response")
	}
	user := data["user"].(map[string]any)
	if user["username"] != username {
		t.Fatalf("expected username %s, got %v", username, user["username"])
	}

	// Login
	loginBody := fmt.Sprintf(`{"email":"%s@test.com","password":"testpass"}`, username)
	w = httptest.NewRecorder()
	r = httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(loginBody))
	r.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatalf("login: expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestIntegrationRefreshToken(t *testing.T) {
	sqlDB := setupIntegrationDB(t)
	defer sqlDB.Close()
	router := newTestRouter(t, sqlDB)

	username := fmt.Sprintf("refresh_%d", time.Now().UnixNano())
	regBody := fmt.Sprintf(`{"username":"%s","email":"%s@test.com","password":"pass"}`, username, username)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(regBody))
	r.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, r)

	var env response.Envelope[any]
	json.Unmarshal(w.Body.Bytes(), &env)
	data := env.Data.(map[string]any)
	refreshToken := data["refresh_token"].(string)

	// Refresh
	refreshBody := fmt.Sprintf(`{"refresh_token":"%s"}`, refreshToken)
	w = httptest.NewRecorder()
	r = httptest.NewRequest("POST", "/api/v1/auth/refresh", strings.NewReader(refreshBody))
	r.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatalf("refresh: expected 200, got %d: %s", w.Code, w.Body.String())
	}

	json.Unmarshal(w.Body.Bytes(), &env)
	if env.Data.(map[string]any)["token"] == nil {
		t.Fatal("expected new token in refresh response")
	}
}

func TestIntegrationUsersMe(t *testing.T) {
	sqlDB := setupIntegrationDB(t)
	defer sqlDB.Close()
	router := newTestRouter(t, sqlDB)

	token := registerAndLogin(t, router, fmt.Sprintf("me_%d", time.Now().UnixNano()))

	w := httptest.NewRecorder()
	r := authRequest("GET", "/api/v1/users/me", token, nil)
	router.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatalf("users/me: expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var env response.Envelope[any]
	json.Unmarshal(w.Body.Bytes(), &env)
	data := env.Data.(map[string]any)
	if data["email"] == nil || data["username"] == nil {
		t.Fatal("expected email and username in response")
	}
}

func TestIntegrationCategories(t *testing.T) {
	sqlDB := setupIntegrationDB(t)
	defer sqlDB.Close()
	router := newTestRouter(t, sqlDB)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/v1/categories", nil)
	router.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatalf("categories: expected 200, got %d", w.Code)
	}

	var env response.Envelope[any]
	json.Unmarshal(w.Body.Bytes(), &env)
	cats, ok := env.Data.([]any)
	if !ok {
		t.Fatal("expected categories array")
	}
	if len(cats) != 12 {
		t.Fatalf("expected 12 categories, got %d", len(cats))
	}
}

func TestIntegrationVideosFeed(t *testing.T) {
	sqlDB := setupIntegrationDB(t)
	defer sqlDB.Close()
	router := newTestRouter(t, sqlDB)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/v1/videos?page=1&size=5", nil)
	router.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatalf("videos: expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var env response.Envelope[any]
	json.Unmarshal(w.Body.Bytes(), &env)
	if env.Pagination == nil {
		t.Fatal("expected pagination in videos response")
	}
}

func TestIntegrationVideoNotFound(t *testing.T) {
	sqlDB := setupIntegrationDB(t)
	defer sqlDB.Close()
	router := newTestRouter(t, sqlDB)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/v1/videos/00000000-0000-0000-0000-000000000000", nil)
	router.ServeHTTP(w, r)

	if w.Code != 404 {
		t.Fatalf("expected 404 for nonexistent video, got %d", w.Code)
	}
}

func TestIntegrationDanmakuSendAndFetch(t *testing.T) {
	sqlDB := setupIntegrationDB(t)
	defer sqlDB.Close()
	router := newTestRouter(t, sqlDB)

	token := registerAndLogin(t, router, fmt.Sprintf("dm_%d", time.Now().UnixNano()))

	// Get an existing video ID
	var videoID string
	sqlDB.QueryRow("SELECT id FROM videos WHERE status = 'published' LIMIT 1").Scan(&videoID)
	if videoID == "" {
		t.Skip("no published video available for danmaku test")
	}

	// Send danmaku
	dmBody := fmt.Sprintf(`{"content":"Integration test danmaku","video_time":1.0,"color":"#FF0000"}`)
	w := httptest.NewRecorder()
	r := authRequest("POST", "/api/v1/videos/"+videoID+"/danmaku", token, strings.NewReader(dmBody))
	router.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatalf("danmaku send: expected 200, got %d: %s", w.Code, w.Body.String())
	}

	// Fetch danmaku
	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", "/api/v1/videos/"+videoID+"/danmaku?t_start=0&t_end=5", nil)
	router.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatalf("danmaku fetch: expected 200, got %d", w.Code)
	}

	var env response.Envelope[any]
	json.Unmarshal(w.Body.Bytes(), &env)
	items, _ := env.Data.([]any)
	if len(items) == 0 {
		t.Fatal("expected at least 1 danmaku")
	}
}

func TestIntegrationComments(t *testing.T) {
	sqlDB := setupIntegrationDB(t)
	defer sqlDB.Close()
	router := newTestRouter(t, sqlDB)

	token := registerAndLogin(t, router, fmt.Sprintf("cmt_%d", time.Now().UnixNano()))

	var videoID string
	sqlDB.QueryRow("SELECT id FROM videos WHERE status = 'published' LIMIT 1").Scan(&videoID)
	if videoID == "" {
		t.Skip("no published video available")
	}

	// Create comment
	body := `{"content":"Integration test comment"}`
	w := httptest.NewRecorder()
	r := authRequest("POST", "/api/v1/videos/"+videoID+"/comments", token, strings.NewReader(body))
	router.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatalf("comment create: expected 200, got %d: %s", w.Code, w.Body.String())
	}

	// List comments
	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", "/api/v1/videos/"+videoID+"/comments", nil)
	router.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatalf("comment list: expected 200, got %d", w.Code)
	}
}

func TestIntegrationSocial(t *testing.T) {
	sqlDB := setupIntegrationDB(t)
	defer sqlDB.Close()
	router := newTestRouter(t, sqlDB)

	token := registerAndLogin(t, router, fmt.Sprintf("soc_%d", time.Now().UnixNano()))

	var videoID string
	sqlDB.QueryRow("SELECT id FROM videos WHERE status = 'published' LIMIT 1").Scan(&videoID)
	if videoID == "" {
		t.Skip("no published video available")
	}

	// Like
	w := httptest.NewRecorder()
	r := authRequest("POST", "/api/v1/videos/"+videoID+"/like", token, nil)
	router.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("like: expected 200, got %d", w.Code)
	}

	// Favorite
	w = httptest.NewRecorder()
	r = authRequest("POST", "/api/v1/videos/"+videoID+"/favorite", token, nil)
	router.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("favorite: expected 200, got %d", w.Code)
	}

	// Watch
	watchBody := `{"progress":15.0,"duration":120}`
	w = httptest.NewRecorder()
	r = authRequest("POST", "/api/v1/videos/"+videoID+"/watch", token, strings.NewReader(watchBody))
	router.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("watch: expected 200, got %d", w.Code)
	}
}

func TestIntegrationSearch(t *testing.T) {
	sqlDB := setupIntegrationDB(t)
	defer sqlDB.Close()
	router := newTestRouter(t, sqlDB)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/v1/search?q=test&type=video", nil)
	router.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatalf("search: expected 200, got %d", w.Code)
	}
}

func TestIntegrationPlaylists(t *testing.T) {
	sqlDB := setupIntegrationDB(t)
	defer sqlDB.Close()
	router := newTestRouter(t, sqlDB)

	token := registerAndLogin(t, router, fmt.Sprintf("pl_%d", time.Now().UnixNano()))

	// Create playlist
	body := `{"name":"Test Playlist","is_public":true}`
	w := httptest.NewRecorder()
	r := authRequest("POST", "/api/v1/playlists", token, strings.NewReader(body))
	router.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatalf("playlist create: expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var env response.Envelope[any]
	json.Unmarshal(w.Body.Bytes(), &env)
	data := env.Data.(map[string]any)
	playlistID := data["id"].(string)

	// Add video to playlist
	var videoID string
	sqlDB.QueryRow("SELECT id FROM videos WHERE status = 'published' LIMIT 1").Scan(&videoID)
	if videoID != "" {
		addBody := fmt.Sprintf(`{"video_id":"%s"}`, videoID)
		w = httptest.NewRecorder()
		r = authRequest("POST", "/api/v1/playlists/"+playlistID+"/videos", token, strings.NewReader(addBody))
		router.ServeHTTP(w, r)
		if w.Code != 200 {
			t.Fatalf("playlist add video: expected 200, got %d: %s", w.Code, w.Body.String())
		}
	}

	// List playlists
	w = httptest.NewRecorder()
	r = authRequest("GET", "/api/v1/playlists", token, nil)
	router.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("playlist list: expected 200, got %d", w.Code)
	}
}

func TestIntegrationAnalytics(t *testing.T) {
	sqlDB := setupIntegrationDB(t)
	defer sqlDB.Close()
	router := newTestRouter(t, sqlDB)

	token := registerAndLogin(t, router, fmt.Sprintf("an_%d", time.Now().UnixNano()))

	w := httptest.NewRecorder()
	r := authRequest("GET", "/api/v1/analytics/overview", token, nil)
	router.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatalf("analytics: expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestIntegrationMetricsEndpoint(t *testing.T) {
	sqlDB := setupIntegrationDB(t)
	defer sqlDB.Close()
	router := newTestRouter(t, sqlDB)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/metrics", nil)
	router.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatalf("metrics: expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "http_requests_total") {
		t.Fatal("expected http_requests_total metric")
	}
}

func TestIntegrationAuthValidation(t *testing.T) {
	sqlDB := setupIntegrationDB(t)
	defer sqlDB.Close()
	router := newTestRouter(t, sqlDB)

	// Register with missing fields
	body := `{"username":"","email":"","password":""}`
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, r)

	if w.Code != 400 {
		t.Fatalf("expected 400 for empty fields, got %d", w.Code)
	}

	// Register with invalid JSON
	w = httptest.NewRecorder()
	r = httptest.NewRequest("POST", "/api/v1/auth/register", strings.NewReader("not json"))
	r.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, r)

	if w.Code != 400 {
		t.Fatalf("expected 400 for invalid JSON, got %d", w.Code)
	}
}

func TestIntegrationCORSHeaders(t *testing.T) {
	sqlDB := setupIntegrationDB(t)
	defer sqlDB.Close()
	router := newTestRouter(t, sqlDB)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/v1/health", nil)
	router.ServeHTTP(w, r)

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatal("expected CORS Allow-Origin header")
	}
}

func TestIntegrationPaginationFormat(t *testing.T) {
	sqlDB := setupIntegrationDB(t)
	defer sqlDB.Close()
	router := newTestRouter(t, sqlDB)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/v1/videos?page=1&size=10", nil)
	router.ServeHTTP(w, r)

	var env response.Envelope[any]
	json.Unmarshal(w.Body.Bytes(), &env)
	p, ok := env.Pagination.(map[string]any)
	if !ok {
		t.Fatal("expected pagination in response")
	}
	if p["page"] != float64(1) || p["size"] != float64(10) {
		t.Fatalf("unexpected pagination: %v", p)
	}
}

func TestIntegrationRefreshInvalidToken(t *testing.T) {
	sqlDB := setupIntegrationDB(t)
	defer sqlDB.Close()
	router := newTestRouter(t, sqlDB)

	body := `{"refresh_token":"invalid-token"}`
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/v1/auth/refresh", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, r)

	if w.Code != 401 {
		t.Fatalf("expected 401 for invalid refresh token, got %d", w.Code)
	}
}

// Test that duplicate registration is rejected
func TestIntegrationDuplicateRegistration(t *testing.T) {
	sqlDB := setupIntegrationDB(t)
	defer sqlDB.Close()
	router := newTestRouter(t, sqlDB)

	username := fmt.Sprintf("dup_%d", time.Now().UnixNano())

	// First registration
	body := fmt.Sprintf(`{"username":"%s","email":"%s@test.com","password":"pass"}`, username, username)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("first register: expected 200, got %d", w.Code)
	}

	// Duplicate registration
	w = httptest.NewRecorder()
	r = httptest.NewRequest("POST", "/api/v1/auth/register", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, r)

	if w.Code != 409 {
		t.Fatalf("expected 409 for duplicate, got %d: %s", w.Code, w.Body.String())
	}
}

// Test that uploading without a file returns error
func TestIntegrationUploadValidation(t *testing.T) {
	sqlDB := setupIntegrationDB(t)
	defer sqlDB.Close()
	router := newTestRouter(t, sqlDB)

	token := registerAndLogin(t, router, fmt.Sprintf("upv_%d", time.Now().UnixNano()))

	// Upload without file
	var buf bytes.Buffer
	body := `--boundary
Content-Disposition: form-data; name="title"

No File Upload
--boundary--`
	buf.WriteString(body)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/v1/videos/upload", &buf)
	r.Header.Set("Authorization", "Bearer "+token)
	r.Header.Set("Content-Type", "multipart/form-data; boundary=boundary")
	router.ServeHTTP(w, r)

	// Should return 400 (file required)
	if w.Code != 400 {
		t.Fatalf("expected 400 for upload without file, got %d: %s", w.Code, w.Body.String())
	}
}
