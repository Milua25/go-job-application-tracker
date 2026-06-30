package middleware

import (
	"log/slog"
	"sync"
	"time"

	"github.com/Milua25/go-job-application-tracker/internal/render"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	mu      sync.RWMutex
	clients map[string]*clientInfo
	limit   int
	reset   time.Duration
}

type clientInfo struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func NewRateLimiter(limit int, reset time.Duration) *RateLimiter {
	r1 := &RateLimiter{
		clients: make(map[string]*clientInfo),
		limit:   limit,
		reset:   reset,
	}
	go r1.cleanupIP()
	return r1
}

// Cleanup removes clients that haven't made requests recently
func (rl *RateLimiter) cleanupIP() {
	ticker := time.NewTicker(rl.reset)
	defer ticker.Stop()
	for {
		<-ticker.C // Wait for the ticker to tick
		rl.mu.Lock()
		for ip, info := range rl.clients {
			if time.Since(info.lastSeen) > rl.reset {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		slog.Info("Client IP", "ip", clientIP)

		rl.mu.Lock()

		if _, exists := rl.clients[clientIP]; !exists {
			// Allow rl.limit requests per second with a burst of rl.limit*2
			rl.clients[clientIP] = &clientInfo{
				limiter:  rate.NewLimiter(rate.Limit(rl.limit), rl.limit*2),
				lastSeen: time.Now(),
			}
		}
		rl.clients[clientIP].lastSeen = time.Now() // Update last seen time for the client
		limiter := rl.clients[clientIP].limiter
		rl.mu.Unlock() // Unlock before checking the limiter

		if !limiter.Allow() {
			render.RateLimitExceededResponse(c, "rate limit exceeded")
			return
		}

		c.Next()
	}
}
