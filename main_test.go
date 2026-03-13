package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/palash287gupta/url-shortner-svc/handler"
	"github.com/palash287gupta/url-shortner-svc/model"
	"github.com/palash287gupta/url-shortner-svc/storage"
	"github.com/palash287gupta/url-shortner-svc/util"
)

func TestGenerateShortCode(t *testing.T) {
	code := util.GenerateShortCode()
	if len(code) != 6 {
		t.Errorf("Expected length 6, got %d", len(code))
	}
}

func TestExtractDomain(t *testing.T) {
	url := "https://www.youtube.com/watch?v=123"
	domain := util.ExtractDomain(url)
	if domain != "www.youtube.com" {
		t.Errorf("Expected www.youtube.com, got %s", domain)
	}
}

func TestShortenHandler(t *testing.T) {
	reqBody := `{"url":"https://www.example.com"}`
	req := httptest.NewRequest("POST", "/shorten", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.ShortenHandler(w, req)

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
	// Clear maps before test
	storage.URLMap = make(map[string]string)
	storage.ReverseMap = make(map[string]string)
	storage.DomainCount = make(map[string]int)

	reqBody := `{"url":"https://www.test.com"}`

	// First request
	req1 := httptest.NewRequest("POST", "/shorten", bytes.NewBufferString(reqBody))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	handler.ShortenHandler(w1, req1)

	var resp1 model.ShortenResponse
	json.NewDecoder(w1.Body).Decode(&resp1)

	// Second request with same URL
	req2 := httptest.NewRequest("POST", "/shorten", bytes.NewBufferString(reqBody))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	handler.ShortenHandler(w2, req2)

	var resp2 model.ShortenResponse
	json.NewDecoder(w2.Body).Decode(&resp2)

	if resp1.ShortURL != resp2.ShortURL {
		t.Errorf("Expected same short URL, got %s and %s", resp1.ShortURL, resp2.ShortURL)
	}
}

func TestRedirectHandler(t *testing.T) {
	// Setup test data
	storage.URLMap = make(map[string]string)
	storage.URLMap["test123"] = "https://www.google.com"

	req := httptest.NewRequest("GET", "/test123", nil)
	w := httptest.NewRecorder()

	handler.RedirectHandler(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("Expected status 302, got %d", w.Code)
	}

	location := w.Header().Get("Location")
	if location != "https://www.google.com" {
		t.Errorf("Expected redirect to https://www.google.com, got %s", location)
	}
}
