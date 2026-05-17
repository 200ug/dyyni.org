package main

import (
	"sync"
	"time"
)

type Limiter struct {
	mu      sync.Mutex
	windows map[string][]time.Time
	limit   int
	window  time.Duration
}

func NewLimiter(limit int, window time.Duration) *Limiter {
	return &Limiter{
		windows: make(map[string][]time.Time),
		limit:   limit,
		window:  window,
	}
}

func (l *Limiter) Allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-l.window)

	entries := l.windows[ip]
	filtered := entries[:0]
	for _, t := range entries {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}

	if len(filtered) >= l.limit {
		l.windows[ip] = filtered
		return false
	}

	filtered = append(filtered, now)
	l.windows[ip] = filtered
	return true
}
