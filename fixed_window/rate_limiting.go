package main

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu        sync.Mutex
	count     int
	limit     int
	window    time.Duration
	resetTime time.Time
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		limit:  limit,
		window: window,
	}
}

func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	if now.After(rl.resetTime) {
		rl.resetTime = now.Add(rl.window)
		rl.count = 0
	}

	if rl.count < rl.limit {
		rl.count++
		return true
	}

	return false
}

type UserLimitStorage struct {
	m      map[string]*RateLimiter
	mu     sync.RWMutex
	limit  int
	window time.Duration
}

func (ul *UserLimitStorage) CheckUser(user string) *RateLimiter {
	ul.mu.RLock()
	rateLimiter, found := ul.m[user]
	ul.mu.RUnlock()

	if found {
		return rateLimiter
	}

	ul.mu.Lock()
	defer ul.mu.Unlock()

	limiter := NewRateLimiter(ul.limit, ul.window)
	ul.m[user] = limiter
	return limiter
}
