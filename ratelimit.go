/*
Simple, thread-safe Go rate-limiter.
Inspired by Antti Huima's algorithm on http://stackoverflow.com/a/668327

Example:

  // Create a new rate-limiter, allowing up-to 10 calls
  // per second
  rl := ratelimit.New(10, time.Second)

  for i:=0; i<20; i++ {
    if rl.Limit() {
      fmt.Println("DOH! Over limit!")
    } else {
      fmt.Println("OK")
    }
  }
*/
package ratelimit

import (
	"sync/atomic"
	"time"
)

// RateLimit instances are thread-safe.
type RateLimiter struct {
	allowance, max, unit, lastCheck int64
}

// New creates a new rate limiter instance
func New(rate int, per time.Duration) *RateLimiter {
	nano := int64(per)
	if nano < 1 {
		nano = int64(time.Second)
	}
	if rate < 1 {
		rate = 1
	}

	return &RateLimiter{
		allowance: int64(rate) * nano, // store our allowance, in ns units
		max:       int64(rate) * nano, // remember our maximum allowance
		unit:      nano,               // remember our unit size

		lastCheck: time.Now().UnixNano(),
	}
}

// Limit returns true if rate was exceeded
func (rl *RateLimiter) Limit() bool {
	// Calculate the number of ns that have passed since our last call
	now := time.Now().UnixNano()
	passed := now - atomic.SwapInt64(&rl.lastCheck, now)

	// Add them to our allowance
	current := atomic.AddInt64(&rl.allowance, passed)

	// Ensure our allowance is not over maximum
	if current > rl.max {
		atomic.AddInt64(&rl.allowance, rl.max-current)
		current = rl.max
	}

	// If our allowance is less than one unit, rate-limit!
	if current < rl.unit {
		return true
	}

	// Not limited, subtract a unit
	atomic.AddInt64(&rl.allowance, -rl.unit)
	return false
}
