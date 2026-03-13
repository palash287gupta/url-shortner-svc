package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var (
	urlMap     = make(map[string]string) // short code -> original URL
	reverseMap = make(map[string]string) // original URL -> short code
	mu         sync.RWMutex
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

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	json.NewDecoder(r.Body).Decode(&req)

	shortCode := generateShortCode()
	urlMap[shortCode] = req.URL

	resp := ShortenResponse{ShortURL: "http://localhost:8080/" + shortCode}
	json.NewEncoder(w).Encode(resp)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/shorten", shortenHandler)

	port := ":8080"
	log.Printf("Starting server on port %s", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
