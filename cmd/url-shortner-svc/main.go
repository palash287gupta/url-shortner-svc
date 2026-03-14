package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := http.ListenAndServe(port, nil); err != nil {
			log.WithField("error", err.Error()).Fatal("Failed to start server")
		}
	}()

	<-quit
	log.WithField("signal", "shutdown").Info("Shutting down server")
}
