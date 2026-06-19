package middlewares

import (
	"context"
	"time"
)

// rateLimiter implements a token bucket algorithm for rate limiting.
type rateLimiter struct {
	tokens     chan struct{}
	refillTime time.Duration
	rate       int
	cancel     context.CancelFunc
}

func newRateLimiter(rate int, refillTime time.Duration) *rateLimiter {
	rl := &rateLimiter{
		tokens:     make(chan struct{}, rate),
		refillTime: refillTime,
		rate:       rate,
	}
	for range rate {
		rl.tokens <- struct{}{}
	}
	ctx, cancel := context.WithCancel(context.Background())
	rl.cancel = cancel
	go rl.startRefill(ctx)
	return rl
}

func (rl *rateLimiter) stop() {
	rl.cancel()
}

// refill the tokens at the specified rate
func (rl *rateLimiter) startRefill(ctx context.Context) {
	ticker := time.NewTicker(rl.refillTime)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for range rl.rate {
				select {
				case rl.tokens <- struct{}{}:
				default:
					// channel is full, stop refilling this tick
				}
			}
		}
	}
}

func (rl *rateLimiter) allow() bool {
	select {
	case <-rl.tokens:
		return true
	default:
		return false
	}
}
