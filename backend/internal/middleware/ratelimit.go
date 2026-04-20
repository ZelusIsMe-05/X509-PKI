package middleware

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter implements simple IP-based rate limiting
type RateLimiter struct {
	mu       sync.Mutex
	attempts map[string][]time.Time
	maxHits  int
	window   time.Duration
}

// NewRateLimiter creates a new rate limiter
// maxHits: maximum number of requests allowed
// window: time window for the rate limit
func NewRateLimiter(maxHits int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		attempts: make(map[string][]time.Time),
		maxHits:  maxHits,
		window:   window,
	}

	// Cleanup expired entries every hour
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		for range ticker.C {
			rl.cleanup()
		}
	}()

	return rl
}

// Allow returns true if the request from ip should be allowed
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// Get attempts for this IP
	attempts := rl.attempts[ip]

	// Remove old attempts outside the window
	validAttempts := []time.Time{}
	for _, attempt := range attempts {
		if attempt.After(cutoff) {
			validAttempts = append(validAttempts, attempt)
		}
	}

	// Check if limit exceeded
	if len(validAttempts) >= rl.maxHits {
		rl.attempts[ip] = validAttempts
		return false
	}

	// Record this attempt
	validAttempts = append(validAttempts, now)
	rl.attempts[ip] = validAttempts

	return true
}

// cleanup removes entries with no recent attempts
func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-2 * rl.window) // Keep entries for 2 windows

	for ip, attempts := range rl.attempts {
		validAttempts := []time.Time{}
		for _, attempt := range attempts {
			if attempt.After(cutoff) {
				validAttempts = append(validAttempts, attempt)
			}
		}
		if len(validAttempts) == 0 {
			delete(rl.attempts, ip)
		} else {
			rl.attempts[ip] = validAttempts
		}
	}
}

// Global rate limiters for auth endpoints
// 5 attempts per 15 minutes
var AuthRateLimiter = NewRateLimiter(5, 15*time.Minute)

// RateLimit middleware checks if the request should be allowed based on IP
func RateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		// Extract IP from X-Forwarded-For if behind proxy
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			ip = xff
		}

		if !AuthRateLimiter.Allow(ip) {
			http.Error(w, "Too many requests. Please try again later.", http.StatusTooManyRequests)
			return
		}

		next(w, r)
	}
}
