package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOK(t *testing.T) {
	w := httptest.NewRecorder()
	OK(w, map[string]string{"key": "value"})

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	ct := w.Header().Get("Content-Type")
	if ct != "application/json; charset=utf-8" {
		t.Fatalf("expected JSON content type, got %q", ct)
	}

	var env Envelope[any]
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("json unmarshal: %v", err)
	}
	if env.Code != 0 {
		t.Fatalf("expected code=0, got %d", env.Code)
	}
	if env.Message != "ok" {
		t.Fatalf("expected message=ok, got %q", env.Message)
	}

	data, ok := env.Data.(map[string]any)
	if !ok {
		t.Fatal("expected data to be a map")
	}
	if data["key"] != "value" {
		t.Fatalf("expected data.key=value, got %v", data["key"])
	}
}

func TestOKNilData(t *testing.T) {
	w := httptest.NewRecorder()
	OK(w, nil)

	var env Envelope[any]
	json.Unmarshal(w.Body.Bytes(), &env)
	if env.Code != 0 {
		t.Fatalf("expected code=0, got %d", env.Code)
	}
}

func TestError(t *testing.T) {
	w := httptest.NewRecorder()
	Error(w, http.StatusNotFound, 40401, "not found")

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}

	var env Envelope[any]
	json.Unmarshal(w.Body.Bytes(), &env)
	if env.Code != 40401 {
		t.Fatalf("expected code=40401, got %d", env.Code)
	}
	if env.Message != "not found" {
		t.Fatalf("expected message='not found', got %q", env.Message)
	}
}

func TestWithPagination(t *testing.T) {
	w := httptest.NewRecorder()
	WithPagination(w, []int{1, 2, 3}, 1, 10, 42)

	var env Envelope[any]
	json.Unmarshal(w.Body.Bytes(), &env)
	if env.Pagination == nil {
		t.Fatal("expected pagination")
	}

	p, ok := env.Pagination.(map[string]any)
	if !ok {
		t.Fatal("expected pagination to be a map")
	}
	if p["page"] != float64(1) || p["size"] != float64(10) || p["total"] != float64(42) {
		t.Fatalf("unexpected pagination: %v", p)
	}
}

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()
	WriteJSON(w, http.StatusCreated, map[string]string{"id": "123"})

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}
