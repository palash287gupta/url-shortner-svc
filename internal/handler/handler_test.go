package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/palash287gupta/url-shortner-svc/internal/model"
	"github.com/palash287gupta/url-shortner-svc/internal/storage"
)

func TestShortenHandler(t *testing.T) {
	h := NewHandler(&Config{BaseURL: "http://localhost:8080"})

	reqBody := `{"url":"https://www.example.com"}`
	req := httptest.NewRequest("POST", "/shorten", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	h.ShortenHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp model.ShortenResponse
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.ShortURL == "" {
		t.Error("Expected short URL, got empty string")
	}
}

func TestDeduplication(t *testing.T) {
	h := NewHandler(&Config{BaseURL: "http://localhost:8080"})

	// Clear maps before test
	storage.URLMap = make(map[string]string)
	storage.ReverseMap = make(map[string]string)
	storage.DomainCount = make(map[string]int)

	reqBody := `{"url":"https://www.test.com"}`

	// First request
	req1 := httptest.NewRequest("POST", "/shorten", bytes.NewBufferString(reqBody))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	h.ShortenHandler(w1, req1)

	var resp1 model.ShortenResponse
	json.NewDecoder(w1.Body).Decode(&resp1)

	// Second request with same URL
	req2 := httptest.NewRequest("POST", "/shorten", bytes.NewBufferString(reqBody))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	h.ShortenHandler(w2, req2)

	var resp2 model.ShortenResponse
	json.NewDecoder(w2.Body).Decode(&resp2)

	if resp1.ShortURL != resp2.ShortURL {
		t.Errorf("Expected same short URL, got %s and %s", resp1.ShortURL, resp2.ShortURL)
	}
}

func TestRedirectHandler(t *testing.T) {
	h := NewHandler(&Config{BaseURL: "http://localhost:8080"})

	// Setup test data
	storage.URLMap = make(map[string]string)
	storage.URLMap["test123"] = "https://www.google.com"

	req := httptest.NewRequest("GET", "/test123", nil)
	w := httptest.NewRecorder()

	h.RedirectHandler(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("Expected status 302, got %d", w.Code)
	}

	location := w.Header().Get("Location")
	if location != "https://www.google.com" {
		t.Errorf("Expected redirect to https://www.google.com, got %s", location)
	}
}
