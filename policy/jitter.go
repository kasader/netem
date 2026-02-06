package policy

import (
	"math/rand/v2"
	"sync/atomic"
	"time"
)

// JitterFunc enables a simple function to satisfy the [Jitter] interface.
type JitterFunc func() time.Duration

// Duration implements the [Jitter] interface.
func (f JitterFunc) Duration() time.Duration { return f() }

// RandomJitter returns a jitter function that selects a random value
// uniformly distributed in the range [-amplitude, +amplitude].
//
// For example, RandomJitter(10*time.Millisecond) will return a duration
// randomly chosen between -10ms and +10ms.
func RandomJitter(amplitude time.Duration) JitterFunc {
	return JitterFunc(func() time.Duration {
		if amplitude == 0 {
			return 0
		}
		n := int64(amplitude)
		delta := rand.Int64N(2 * n)
		return time.Duration(delta - n)
	})
}

// JitterVar is a thread-safe, mutable [Jitter] provider.
// It allows you to change the jitter of a running simulation.
//
// Uses the [RandomJitter] policy. For other policies, please implement a
// custom LossVar implementation.
type JitterVar struct{ val atomic.Int64 }

// Set updates the jitter safely.
func (v *JitterVar) Set(d time.Duration) { v.val.Store(int64(d)) }

// Duration implements the [Latency] interface.
func (v *JitterVar) Duration() time.Duration {
	amplitude := time.Duration(v.val.Load())
	return RandomJitter(amplitude)()
}
