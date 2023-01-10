package ratelimit_test

import (
	"testing"
	"time"

	. "github.com/bsm/ratelimit"
)

func delta(x, y int) int {
	if x > y {
		return x - y
	}
	return y - x
}

func TestRateLimiter_smallRates(t *testing.T) {
	rl := New(10, time.Minute)
	count := 0
	for !rl.Limit() {
		count++
	}
	if exp, got := 10, count; exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}
}

func TestRateLimiter_largeRates(t *testing.T) {
	rl := New(100_000, time.Hour)
	count := 0
	for !rl.Limit() {
		count++
	}
	if exp, got := 100_000, count; delta(exp, got) > 10 {
		t.Errorf("expected %v, got %v", exp, got)
	}
}

func TestRateLimiter_largeIntevals(t *testing.T) {
	rl := New(10, 360*24*time.Hour)
	count := 0
	for !rl.Limit() {
		count++
	}
	if exp, got := 10, count; exp != got {
		t.Errorf("expected %v, got %v", exp, got)
	}
}

func TestRateLimiter_increaseAllowance(t *testing.T) {
	n := 25
	rl := New(n, 50*time.Millisecond)
	for i := 0; i < n; i++ {
		if rl.Limit() {
			t.Errorf("expected no limit on cycle %d", i+1)
		}
	}

	if !rl.Limit() {
		t.Errorf("expected limit on cycle %d", n+1)
	}

	time.Sleep(20 * time.Millisecond)
	if rl.Limit() {
		t.Errorf("expected no limit after delay")
	}
}

func TestRateLimiter_spreadAllowance(t *testing.T) {
	rl := New(5, 10*time.Millisecond)
	start := time.Now()
	count := 0
	for time.Since(start) < 100*time.Millisecond {
		if !rl.Limit() {
			count++
		}
	}
	if exp, got := 54, count; delta(exp, got) > 1 {
		t.Errorf("expected %v, got %v", exp, got)
	}
}

func TestRateLimiter_Undo(t *testing.T) {
	n := 5
	rl := New(n, time.Minute)
	for i := 0; i < n; i++ {
		if rl.Limit() {
			t.Errorf("expected no limit on cycle %d", i+1)
		}
	}
	if !rl.Limit() {
		t.Errorf("expected limit on cycle %d", n+1)
	}

	rl.Undo()
	if rl.Limit() {
		t.Error("expected no limit after undo")
	}
	if !rl.Limit() {
		t.Error("expected to limit again")
	}
}

func BenchmarkLimit(b *testing.B) {
	rl := New(1000, time.Second)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rl.Limit()
	}
}
