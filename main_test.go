package main

import "testing"

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
