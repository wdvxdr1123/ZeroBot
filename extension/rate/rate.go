// Package rate provides a rate limiter hooker,
// this package is based on golang.org/x/time/rate
package rate

import (
	"sync"
	"time"

	"github.com/wdvxdr1123/ZeroBot/extension/ttl"
)

// LimiterManager ...
type LimiterManager[K comparable] struct {
	limiters *ttl.Cache[K, *Limiter]
	interval time.Duration
	burst    int
}

// NewManager ..
func NewManager[K comparable](interval time.Duration, burst int) *LimiterManager[K] {
	return &LimiterManager[K]{
		limiters: ttl.NewCache[K, *Limiter](interval * time.Duration(burst)),
		interval: interval,
		burst:    burst,
	}
}

// Load ...
func (l *LimiterManager[K]) Load(key K) *Limiter {
	if val := l.limiters.Get(key); val != nil {
		return val
	}
	val := NewLimiter(l.interval, l.burst)
	l.limiters.Set(key, val)
	return val
}

// Limiter controls the frequency of handling events.
type Limiter struct {
	sync.Mutex
	limit    float64
	tokens   float64
	burst    int
	lastTime time.Time
}

// Tokens returns the left token of limiter.
func (lim *Limiter) Tokens() float64 {
	return lim.tokens
}

// NewLimiter returns a new Limiter Pointer.
func NewLimiter(interval time.Duration, burst int) *Limiter {
	return &Limiter{
		limit:    every(interval),
		burst:    burst,
		tokens:   float64(burst),
		lastTime: time.Now(),
	}
}

// Acquire ...
func (lim *Limiter) Acquire() bool {
	return lim.AcquireN(1)
}

// AcquireN ...
func (lim *Limiter) AcquireN(n int) bool {
	lim.Lock()
	defer lim.Unlock()
	lim.advance(time.Now())
	nf := float64(n)
	if lim.tokens >= nf {
		lim.tokens -= nf
		return true
	}
	return false
}

func (lim *Limiter) advance(now time.Time) {
	last := lim.lastTime
	elapsed := now.Sub(last)
	if maxElapsed := lim.durationFromTokens(float64(lim.burst) - lim.tokens); elapsed > maxElapsed {
		elapsed = maxElapsed
	}
	delta := lim.tokensFromDuration(elapsed)
	tokens := lim.tokens + delta
	if burst := float64(lim.burst); tokens > burst {
		tokens = burst
	}
	lim.tokens = tokens
	lim.lastTime = now
}

func every(interval time.Duration) float64 {
	return 1 / interval.Seconds()
}

func (lim *Limiter) durationFromTokens(tokens float64) time.Duration {
	seconds := tokens / lim.limit
	return time.Nanosecond * time.Duration(1e9*seconds)
}

func (lim *Limiter) tokensFromDuration(d time.Duration) float64 {
	sec := float64(d/time.Second) * lim.limit
	nSec := float64(d%time.Second) * lim.limit
	return sec + nSec/1e9
}
