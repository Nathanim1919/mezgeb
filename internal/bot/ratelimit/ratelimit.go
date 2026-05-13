package ratelimit

import (
	"sync"
	"time"
)

// Limiter is a simple per-user rate limiter using a sliding window.
type Limiter struct {
	mu       sync.Mutex
	requests map[int64][]time.Time
	max      int
	window   time.Duration
}

func New(max int, window time.Duration) *Limiter {
	return &Limiter{
		requests: make(map[int64][]time.Time),
		max:      max,
		window:   window,
	}
}

// Allow returns true if the user hasn't exceeded the rate limit.
func (l *Limiter) Allow(userID int64) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-l.window)

	// Remove expired entries
	times := l.requests[userID]
	valid := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= l.max {
		l.requests[userID] = valid
		return false
	}

	l.requests[userID] = append(valid, now)
	return true
}
