package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenerateShortCode(t *testing.T) {
	code := generateShortCode()
	if len(code) != 6 {
		t.Errorf("Expected length 6, got %d", len(code))
	}
}

func TestExtractDomain(t *testing.T) {
	url := "https://www.youtube.com/watch?v=123"
	domain := extractDomain(url)
	if domain != "www.youtube.com" {
		t.Errorf("Expected www.youtube.com, got %s", domain)
	}
}

func TestShortenHandler(t *testing.T) {
	reqBody := `{"url":"https://www.example.com"}`
	req := httptest.NewRequest("POST", "/shorten", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	shortenHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp ShortenResponse
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.ShortURL == "" {
		t.Error("Expected short URL, got empty string")
	}
}
