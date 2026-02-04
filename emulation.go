package netem

import (
	"math"
	"math/rand/v2"
	"sync/atomic"
	"time"
)

// --- Bandwidth

// Bandwidth models the capacity of the link.
type Bandwidth interface {
	// Limit returns the allowed throughput in bits per second.
	Limit() uint64
}

// BandwidthFunc enables a simple function to satisfy the [Bandwidth] interface.
type BandwidthFunc func() uint64

func (f BandwidthFunc) Limit() uint64 { return f() }

// StaticBandwidth returns a constant throughput.
func StaticBandwidth(bps uint64) Bandwidth { return BandwidthFunc(func() uint64 { return bps }) }

// BandwidthVar is a thread-safe, mutable [Bandwidth] provider.
// It allows you to change the bandwidth of a running simulation.
type BandwidthVar struct{ val atomic.Uint64 }

// Set updates the bandwidth safely.
func (v *BandwidthVar) Set(bandwidth uint64) { v.val.Store(bandwidth) }

// Limit implements the [Bandwidth] interface.
func (v *BandwidthVar) Limit() uint64 { return v.val.Load() }

// --- Latency

// Latency models the delay of a network transmission.
type Latency interface {
	// Duration returns the delay for the current operation.
	Duration() time.Duration
}

// TODO: Add sine-wave latency func, etc. more utilities!

// LatencyFunc enables a simple function to satisfy the [Latency] interface.
type LatencyFunc func() time.Duration

func (f LatencyFunc) Duration() time.Duration { return f() }

// StaticLatency returns a constant delay.
func StaticLatency(d time.Duration) Latency { return LatencyFunc(func() time.Duration { return d }) }

// LatencyVar is a thread-safe, mutable [Latency] provider.
// It allows you to change the latency of a running simulation.
type LatencyVar struct{ val atomic.Int64 }

// Set updates the latency safely.
func (v *LatencyVar) Set(latency time.Duration) { v.val.Store(int64(latency)) }

// Duration implements the [Latency] interface.
func (v *LatencyVar) Duration() time.Duration { return time.Duration(v.val.Load()) }

// --- Jitter

// Jitter models the variance in transmission delay.
type Jitter interface {
	// Duration returns the random variance to add to the latency.
	Duration() time.Duration
}

// JitterFunc enables a simple function to satisfy the [Jitter] interface.
type JitterFunc func() time.Duration

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

// --- Loss

// Loss models the unreliability of a datagram link.
type Loss interface {
	// Drop returns true if the current datagram should be discarded.
	Drop() bool
}

// LossFunc enables a simple function to satisfy the [Loss] interface.
type LossFunc func() bool

func (f LossFunc) Drop() bool { return f() }

// RandomLoss returns a function that drops datagrams with probability rate (0.0 to 1.0).
func RandomLoss(rate float64) LossFunc {
	return LossFunc(func() bool {
		return rand.Float64() < rate
	})
}

// LossVar is a thread-safe, mutable [Loss] provider.
// It allows you to change the random loss rate of a running simulation.
//
// Uses the [RandomLoss] policy. For other policies, please implement a
// custom LossVar implementation.
type LossVar struct {
	// We store our loss rate (f64) within an [atomic.Uint64].
	// see: https://github.com/golang/go/issues/21996
	val atomic.Uint64
}

// Set updates the loss rate safely.
func (v *LossVar) Set(rate float64) { v.val.Store(math.Float64bits(rate)) }

// Duration implements the [Loss] interface.
func (v *LossVar) Drop() bool {
	rate := math.Float64frombits(v.val.Load())
	return RandomLoss(rate)()
}

// --- Fault

// Fault models the stability of a connection.
type Fault interface {
	// ShouldClose returns true if the connection should be severed abruptly.
	ShouldClose() bool
}

// FaultFunc enables a simple function to satisfy the [Fault] interface.
type FaultFunc func() bool

func (f FaultFunc) ShouldClose() bool { return f() }

// RandomFault returns a function that closes connections with probability rate (0.0 to 1.0).
func RandomClose(rate float64) FaultFunc {
	return func() bool {
		return rand.Float64() < rate
	}
}

// FaultVar is a thread-safe, mutable [Fault] provider.
// It allows you to change the random fault rate of a running simulation.
//
// Uses the [RandomClose] policy. For other policies, please implement a
// custom FaultVar implementation.
type FaultVar struct {
	val atomic.Uint64
}

// Set updates the fault rate safely.
func (v *FaultVar) Set(rate float64) { v.val.Store(math.Float64bits(rate)) }

// Duration implements the [Fault] interface.
func (v *FaultVar) Drop() bool {
	rate := math.Float64frombits(v.val.Load())
	return RandomLoss(rate)()
}
