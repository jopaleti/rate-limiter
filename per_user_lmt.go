package main

import "time"

// Example with map of rate limiters per user
type RateLimiterManager struct {
    limiters map[string]*TokenBucketRateLimiter
}

func (m *RateLimiterManager) GetLimiter(userID string) *TokenBucketRateLimiter {
    limiter, exists := m.limiters[userID]
    if !exists {
        limiter = NewTokenBucketRateLimiter(100, time.Minute/100) // Default 100 requests/minute
        m.limiters[userID] = limiter
    }
    return limiter
}