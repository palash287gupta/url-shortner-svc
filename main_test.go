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

func TestDeduplication(t *testing.T) {
	// Clear maps before test
	urlMap = make(map[string]string)
	reverseMap = make(map[string]string)
	domainCount = make(map[string]int)

	reqBody := `{"url":"https://www.test.com"}`

	// First request
	req1 := httptest.NewRequest("POST", "/shorten", bytes.NewBufferString(reqBody))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	shortenHandler(w1, req1)

	var resp1 ShortenResponse
	json.NewDecoder(w1.Body).Decode(&resp1)

	// Second request with same URL
	req2 := httptest.NewRequest("POST", "/shorten", bytes.NewBufferString(reqBody))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	shortenHandler(w2, req2)

	var resp2 ShortenResponse
	json.NewDecoder(w2.Body).Decode(&resp2)

	if resp1.ShortURL != resp2.ShortURL {
		t.Errorf("Expected same short URL, got %s and %s", resp1.ShortURL, resp2.ShortURL)
	}
}
