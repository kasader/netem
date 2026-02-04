package netem

import (
	"math"
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

// StaticJitter returns a constant jitter variance.
func StaticJitter(d time.Duration) Jitter { return JitterFunc(func() time.Duration { return d }) }

// JitterVar is a thread-safe, mutable [Jitter] provider.
// It allows you to change the jitter of a running simulation.
type JitterVar struct{ val atomic.Int64 }

// Set updates the jitter safely.
func (v *JitterVar) Set(d time.Duration) { v.val.Store(int64(d)) }

// Duration implements the [Latency] interface.
func (v *JitterVar) Duration() time.Duration { return time.Duration(v.val.Load()) }

// --- Loss

// Loss models the unreliability of a datagram link.
type Loss interface {
	// Drop returns true if the current datagram should be discarded.
	Drop() bool
}

// LossFunc enables a simple function to satisfy the [Loss] interface.
type LossFunc func() bool

func (f LossFunc) Drop() bool { return f() }

// StaticLoss returns a datagram drop decision via a static loss rate.
//
// TODO:
func StaticLoss(f float64) Loss {
	return LossFunc(func() bool {
		panic("not implemented yet")
		// return d
	})
}

// LossVar is a thread-safe, mutable [Loss] provider.
// It allows you to change the loss rate of a running simulation.
type LossVar struct {
	// We store our loss rate (f64) within an [atomic.Uint64].
	// see: https://github.com/golang/go/issues/21996
	val atomic.Uint64
}

// Set updates the loss rate safely.
func (v *LossVar) Set(loss float64) { v.val.Store(math.Float64bits(loss)) }

// Duration implements the [Latency] interface.
func (v *LossVar) Drop() bool {
	panic("not implemented yet")
	// math.Float64frombits(v.val.Load())
}

// --- Fault

// Fault models the stability of a connection.
type Fault interface {
	// ShouldClose returns true if the connection should be severed abruptly.
	ShouldClose() bool
}
