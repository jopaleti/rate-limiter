package main

import (
	"fmt"
	"sync"
	"time"
)

// TokenBucketRateLimiter struct to handle rate limiting
type TokenBucketRateLimiter struct {
	capacity int // Maximum number of tokens
	tokens int // Current number of tokens
	refillRate time.Duration // Time interval to add one token
	lastRefill time.Time // Last time tokens were refilled
	mutex sync.Mutex // For thread-safe access
}

// NewTokenBucketRateLimiter creates a new rate limiter with a given capacity and refill rate
func NewTokenBucketRateLimiter(capacity int, refillRate time.Duration) *TokenBucketRateLimiter {
	return &TokenBucketRateLimiter{
		capacity: capacity,
		tokens: capacity, // start with a full bucket
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Refill adds tokens based on elapsed time since last refill
func (rl *TokenBucketRateLimiter) Refill() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)

	// Calculate how many tokens to add based on the elapsed time
	newTokens := int(elapsed / rl.refillRate)
	if newTokens > 0 {
		rl.tokens = min(rl.capacity, rl.tokens+newTokens) // Ensure we don't exceed capacity
		rl.lastRefill = time.Now() // Reset the last refill rate
	}
}

// Allow checks if a request can proceed (i.e., if there's a token available)
func (rl *TokenBucketRateLimiter) Allow() bool {
	rl.Refill() // Refill tokens based on time

	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	if rl.tokens > 0 {
		rl.tokens-- // consume a token
		return true // Request is allowed
	}

	return false // No token available, request is rate-limited
}

// Rate Limiter in Action
func main() {
	// Create a rate limiter allowing 5 requests per second
	rateLimiter := NewTokenBucketRateLimiter(5, time.Second/5)

	// Simulate incoming requests
	for i := 0; i < 10; i++ {
        if rateLimiter.Allow() {
            fmt.Printf("Request %d allowed\n", i+1)
        } else {
            fmt.Printf("Request %d rate-limited\n", i+1)
        }
        time.Sleep(200 * time.Millisecond) // Simulate a delay between requests
    }
}