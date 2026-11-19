package router

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type visitor struct {
	tokens   float64
	lastSeen time.Time
}

type rateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     float64 // tokens restored per second
	burst    int     // max tokens (bucket size)
	now      func() time.Time
}

func newRateLimiter(requestsPerMinute int, burst int) *rateLimiter {
	rl := &rateLimiter{
		visitors: make(map[string]*visitor),
		rate:     float64(requestsPerMinute) / 60.0,
		burst:    burst,
		now:      time.Now,
	}
	go rl.cleanup()
	return rl
}

// allow checks whether the given IP has tokens available.
func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := rl.now()
	v, exists := rl.visitors[ip]
	if !exists {
		rl.visitors[ip] = &visitor{tokens: float64(rl.burst) - 1, lastSeen: now}
		return true
	}

	elapsed := now.Sub(v.lastSeen).Seconds()
	v.lastSeen = now
	v.tokens += elapsed * rl.rate
	if v.tokens > float64(rl.burst) {
		v.tokens = float64(rl.burst)
	}

	if v.tokens < 1 {
		return false
	}
	v.tokens--
	return true
}

// cleanup removes stale entries every 5 minutes to prevent memory leaks.
func (rl *rateLimiter) cleanup() {
	for {
		time.Sleep(5 * time.Minute)
		rl.evictStale(10 * time.Minute)
	}
}

// evictStale drops visitors not seen for longer than maxIdle.
func (rl *rateLimiter) evictStale(maxIdle time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := rl.now()
	for ip, v := range rl.visitors {
		if now.Sub(v.lastSeen) > maxIdle {
			delete(rl.visitors, ip)
		}
	}
}

// RateLimit returns a Gin middleware that enforces the given rate limit per IP.
func RateLimit(requestsPerMinute int, burst int) gin.HandlerFunc {
	limiter := newRateLimiter(requestsPerMinute, burst)
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.allow(ip) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests, try again later",
			})
			return
		}
		c.Next()
	}
}
