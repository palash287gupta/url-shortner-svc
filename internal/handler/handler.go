package handler

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/palash287gupta/url-shortner-svc/internal/model"
	"github.com/palash287gupta/url-shortner-svc/internal/storage"
	"github.com/palash287gupta/url-shortner-svc/internal/util"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	BaseURL string
}

type Handler struct {
	config *Config
}

func NewHandler(cfg *Config) *Handler {
	return &Handler{
		config: cfg,
	}
}

func (h *Handler) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.WithFields(log.Fields{
			"method": r.Method,
		}).Warn("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Failed to decode request")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		log.Warn("Empty URL provided")
		http.Error(w, "URL cannot be empty", http.StatusBadRequest)
		return
	}

	storage.Mu.RLock()
	existingShortCode, exists := storage.ReverseMap[req.URL]
	storage.Mu.RUnlock()

	if exists {
		log.WithFields(log.Fields{
			"url":       req.URL,
			"shortCode": existingShortCode,
		}).Info("URL already shortened, returning existing short code")
		resp := model.ShortenResponse{
			ShortURL: h.config.BaseURL + "/" + existingShortCode,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	shortCode := util.GenerateShortCode()

	storage.Mu.Lock()
	storage.URLMap[shortCode] = req.URL
	storage.ReverseMap[req.URL] = shortCode
	domain := util.ExtractDomain(req.URL)
	storage.DomainCount[domain]++
	storage.Mu.Unlock()

	log.WithFields(log.Fields{
		"url":       req.URL,
		"shortCode": shortCode,
		"domain":    domain,
	}).Info("URL shortened successfully")

	resp := model.ShortenResponse{
		ShortURL: h.config.BaseURL + "/" + shortCode,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	shortCode := strings.TrimPrefix(r.URL.Path, "/")

	if shortCode == "" {
		log.Warn("Empty short code provided")
		http.Error(w, "Short code not found", http.StatusNotFound)
		return
	}

	storage.Mu.RLock()
	originalURL, exists := storage.URLMap[shortCode]
	storage.Mu.RUnlock()

	if !exists {
		log.WithFields(log.Fields{
			"shortCode": shortCode,
		}).Warn("Short code not found")
		http.Error(w, "Short code not found", http.StatusNotFound)
		return
	}

	log.WithFields(log.Fields{
		"shortCode":   shortCode,
		"originalURL": originalURL,
	}).Info("Redirecting to original URL")

	http.Redirect(w, r, originalURL, http.StatusFound)
}

func (h *Handler) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	storage.Mu.RLock()
	defer storage.Mu.RUnlock()

	type domainMetric struct {
		domain string
		count  int
	}

	var domains []domainMetric
	for domain, count := range storage.DomainCount {
		domains = append(domains, domainMetric{domain: domain, count: count})
	}

	sort.Slice(domains, func(i, j int) bool {
		return domains[i].count > domains[j].count
	})

	log.WithFields(log.Fields{
		"totalDomains": len(domains),
	}).Info("Metrics requested")

	w.Header().Set("Content-Type", "text/plain")
	for i := 0; i < len(domains) && i < 3; i++ {
		w.Write([]byte(domains[i].domain + ": " + strconv.Itoa(domains[i].count) + "\n"))
	}
}
