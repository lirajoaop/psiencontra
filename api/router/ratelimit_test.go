package router

import (
	"testing"
	"time"
)

// fakeClock lets tests advance time deterministically.
type fakeClock struct{ t time.Time }

func (c *fakeClock) now() time.Time          { return c.t }
func (c *fakeClock) advance(d time.Duration) { c.t = c.t.Add(d) }

func newTestLimiter(perMinute, burst int) (*rateLimiter, *fakeClock) {
	clock := &fakeClock{t: time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)}
	rl := &rateLimiter{
		visitors: map[string]*visitor{},
		rate:     float64(perMinute) / 60.0,
		burst:    burst,
		now:      clock.now,
	}
	return rl, clock
}

func TestRateLimiter_Allow(t *testing.T) {
	type step struct {
		advance time.Duration
		ip      string
		want    bool
	}
	const ip = "203.0.113.1"

	tests := []struct {
		name      string
		perMinute int
		burst     int
		steps     []step
	}{
		{
			name:      "burst is allowed then blocked",
			perMinute: 60, burst: 3,
			steps: []step{{0, ip, true}, {0, ip, true}, {0, ip, true}, {0, ip, false}},
		},
		{
			name:      "one token refills after 1/rate seconds",
			perMinute: 60, burst: 1,
			steps: []step{
				{0, ip, true},
				{500 * time.Millisecond, ip, false}, // half a token
				{600 * time.Millisecond, ip, true},  // a full token again
				{0, ip, false},
			},
		},
		{
			name:      "auth scope: 10/min, burst 5, one token every 6s",
			perMinute: 10, burst: 5,
			steps: []step{
				{0, ip, true}, {0, ip, true}, {0, ip, true}, {0, ip, true}, {0, ip, true},
				{0, ip, false},
				{5 * time.Second, ip, false}, // 5/6 of a token
				{2 * time.Second, ip, true},  // past 6s: one token back
				{0, ip, false},
			},
		},
		{
			name:      "answer scope: 5/min, burst 2, one token every 12s",
			perMinute: 5, burst: 2,
			steps: []step{
				{0, ip, true}, {0, ip, true}, {0, ip, false},
				{11 * time.Second, ip, false},
				{2 * time.Second, ip, true},
			},
		},
		{
			name:      "bucket never exceeds burst after long idle",
			perMinute: 60, burst: 2,
			steps: []step{
				{0, ip, true},
				{time.Hour, ip, true},
				{0, ip, true},
				{0, ip, false},
			},
		},
		{
			name:      "different IPs have independent buckets",
			perMinute: 60, burst: 1,
			steps: []step{
				{0, ip, true},
				{0, ip, false},
				{0, "198.51.100.7", true},
				{0, "198.51.100.7", false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rl, clock := newTestLimiter(tt.perMinute, tt.burst)
			for i, s := range tt.steps {
				clock.advance(s.advance)
				if got := rl.allow(s.ip); got != s.want {
					t.Fatalf("step %d (%s after +%v): allow = %v, want %v", i, s.ip, s.advance, got, s.want)
				}
			}
		})
	}
}

func TestRateLimiter_EvictStale(t *testing.T) {
	rl, clock := newTestLimiter(60, 1)
	rl.allow("old")
	clock.advance(11 * time.Minute)
	rl.allow("fresh")

	rl.evictStale(10 * time.Minute)

	if _, ok := rl.visitors["old"]; ok {
		t.Error("stale visitor was not evicted")
	}
	if _, ok := rl.visitors["fresh"]; !ok {
		t.Error("fresh visitor was evicted")
	}
}
