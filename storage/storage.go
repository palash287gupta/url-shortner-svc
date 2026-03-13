package storage

import "sync"

var (
	URLMap      = make(map[string]string) // short code -> original URL
	ReverseMap  = make(map[string]string) // original URL -> short code
	DomainCount = make(map[string]int)    // domain -> count
	Mu          sync.RWMutex
)
