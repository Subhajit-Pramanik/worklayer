package middleware

import (
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"github.com/vyolayer/vyolayer/pkg/logger"
	"github.com/vyolayer/vyolayer/pkg/response"
)

// RateLimiter is a simple in-memory, sliding-window rate limiter.
// It is safe for concurrent use.
type RateLimiter struct {
	mu         sync.Mutex
	requests   map[string][]time.Time
	max        int
	windowSize time.Duration
	key        string
	logger     logger.AppLogger
}

// NewRateLimiter returns a configured RateLimiter.
// max      – maximum number of requests allowed within window.
// window   – duration of the sliding window (e.g. 1*time.Minute).
func NewRateLimiter(max int, window time.Duration, key string) *RateLimiter {
	rl := &RateLimiter{
		requests:   make(map[string][]time.Time),
		max:        max,
		windowSize: window,
		logger:     *logger.ForService("RateLimiter"),
	}

	if key != "" {
		rl.key = key
	}

	// Background goroutine to prune stale entries and avoid unbounded growth.
	go func() {
		ticker := time.NewTicker(window)
		defer ticker.Stop()
		for range ticker.C {
			rl.prune()
		}
	}()

	return rl
}

func (rl *RateLimiter) prune() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	cutoff := time.Now().Add(-rl.windowSize)
	for key, timestamps := range rl.requests {
		var fresh []time.Time
		for _, ts := range timestamps {
			if ts.After(cutoff) {
				fresh = append(fresh, ts)
			}
		}
		if len(fresh) == 0 {
			delete(rl.requests, key)
		} else {
			rl.requests[key] = fresh
		}
	}
}

// Handler returns a Fiber middleware function.
// key defaults to the caller's IP address.
func (rl *RateLimiter) Handler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := c.IP()
		if rl.key != "" {
			key = key + ":" + strings.ToLower(rl.key)
		}

		rl.mu.Lock()
		now := time.Now()
		cutoff := now.Add(-rl.windowSize)

		// Slide the window: keep only requests within the window.
		var fresh []time.Time
		for _, ts := range rl.requests[key] {
			if ts.After(cutoff) {
				fresh = append(fresh, ts)
			}
		}

		rl.logger.Debug("Rate limiter", map[string]interface{}{
			"key":      key,
			"max":      rl.max,
			"window":   rl.windowSize,
			"requests": len(fresh),
		})

		if len(fresh) >= rl.max {
			rl.mu.Unlock()
			return response.Error(c, errors.TooManyRequests("too many requests, please slow down"))
		}

		rl.requests[key] = append(fresh, now)
		rl.mu.Unlock()

		return c.Next()
	}
}
