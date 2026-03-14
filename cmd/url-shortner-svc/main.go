package main

import (
	"net/http"

	"github.com/palash287gupta/url-shortner-svc/internal/handler"

	log "github.com/sirupsen/logrus"
)

func main() {
	cfg := LoadConfig()

	h := handler.NewHandler(&handler.Config{
		BaseURL: cfg.BaseURL,
	})

	http.HandleFunc("/shorten", h.ShortenHandler)
	http.HandleFunc("/metrics", h.MetricsHandler)
	http.HandleFunc("/", h.RedirectHandler)

	port := ":" + cfg.Port
	log.WithField("port", port).Info("Starting URL Shortener Service")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
