package main

import (
	"net/http"

	"github.com/palash287gupta/url-shortner-svc/config"
	"github.com/palash287gupta/url-shortner-svc/handler"

	log "github.com/sirupsen/logrus"
)

func main() {
	cfg := config.LoadConfig()

	http.HandleFunc("/shorten", handler.ShortenHandler)
	http.HandleFunc("/metrics", handler.MetricsHandler)
	http.HandleFunc("/", handler.RedirectHandler)

	port := ":" + cfg.Port
	log.WithField("port", port).Info("Starting URL Shortener Service")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
