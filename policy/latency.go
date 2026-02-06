package policy

import (
	"sync/atomic"
	"time"
)

// TODO: Add sine-wave latency func, etc. more utilities!

// LatencyFunc enables a simple function to satisfy the [Latency] interface.
type LatencyFunc func() time.Duration

// Duration implements the [Latency] interface.
func (f LatencyFunc) Duration() time.Duration { return f() }

// StaticLatency returns a constant delay.
func StaticLatency(d time.Duration) LatencyFunc {
	return LatencyFunc(func() time.Duration { return d })
}

// LatencyVar is a thread-safe, mutable [Latency] provider.
// It allows you to change the latency of a running simulation.
type LatencyVar struct{ val atomic.Int64 }

// Set updates the latency safely.
func (v *LatencyVar) Set(latency time.Duration) { v.val.Store(int64(latency)) }

// Duration implements the [Latency] interface.
func (v *LatencyVar) Duration() time.Duration { return time.Duration(v.val.Load()) }
