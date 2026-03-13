package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	urlMap      = make(map[string]string) // short code -> original URL
	reverseMap  = make(map[string]string) // original URL -> short code
	domainCount = make(map[string]int)    // domain -> count
	mu          sync.RWMutex
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

func generateShortCode() string {
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func extractDomain(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) > 2 {
		return parts[2]
	}
	return url
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// Check if URL already shortened
	if existingCode, found := reverseMap[req.URL]; found {
		resp := ShortenResponse{ShortURL: "http://localhost:8080/" + existingCode}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	shortCode := generateShortCode()
	urlMap[shortCode] = req.URL
	reverseMap[req.URL] = shortCode

	// Track domain count
	domain := extractDomain(req.URL)
	domainCount[domain]++

	resp := ShortenResponse{ShortURL: "http://localhost:8080/" + shortCode}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	shortCode := strings.TrimPrefix(path, "/")

	if shortCode == "" {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	mu.RLock()
	originalURL, found := urlMap[shortCode]
	mu.RUnlock()

	if !found {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	// Just print all domains and their counts for now
	w.Header().Set("Content-Type", "text/plain")
	for domain, count := range domainCount {
		w.Write([]byte(domain + ": " + string(rune(count)) + "\n"))
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/metrics", metricsHandler)
	http.HandleFunc("/", redirectHandler)

	port := ":8080"
	log.Printf("Starting server on port %s", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
