package ratelimiter

import (
	"context"
	"sync"
	"time"
)

type RateLimiter struct {
	requests int
	window   time.Duration
	mu       sync.Mutex
	buckets  map[string]*bucket
}

type bucket struct {
	tokens    int
	lastReset time.Time
}

func NewRateLimiter(requests int, windowSeconds int) *RateLimiter {
	return &RateLimiter{
		requests: requests,
		window:   time.Duration(windowSeconds) * time.Second,
		buckets:  make(map[string]*bucket),
	}
}

func (r *RateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	b, exists := r.buckets[key]

	if !exists || now.Sub(b.lastReset) >= r.window {
		r.buckets[key] = &bucket{
			tokens:    r.requests - 1,
			lastReset: now,
		}
		return true, nil
	}

	if b.tokens > 0 {
		b.tokens--
		return true, nil
	}

	return false, nil
}