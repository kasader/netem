package netem

import (
	"math"
	"sync/atomic"
	"time"
)

// Bandwidth models the capacity of the link.
type Bandwidth interface {
	// Limit returns the allowed throughput in bits per second.
	Limit() uint64
}

// Latency models the delay of a network transmission.
type Latency interface {
	// Duration returns the delay for the current operation.
	Duration() time.Duration
}

// Jitter models the variance in transmission delay.
type Jitter interface {
	// Duration returns the random variance to add to the latency.
	Duration() time.Duration
}

// Loss models the unreliability of a datagram link.
type Loss interface {
	// Drop returns true if the current datagram should be discarded.
	Drop() bool
}

// Fault models the stability of a connection.
type Fault interface {
	// ShouldClose returns true if the connection should be severed abruptly.
	ShouldClose() bool
}

// LatencyFunc enables a simple function to satisfy LatencyGenerator.
type LatencyFunc func() time.Duration

func (f LatencyFunc) Generate() time.Duration {
	return f()
}

type BandwidthVar struct {
	val atomic.Uint64
}

func (v *BandwidthVar) Set(bandwidth uint64) { v.val.Store(bandwidth) }
func (v *BandwidthVar) Bandwidth() uint64    { return v.val.Load() }

// LatencyVar is a thread-safe, mutable LatencyGenerator.
// It allows you to change the latency of a running simulation.
type LatencyVar struct {
	val atomic.Int64
}

// Set updates the latency safely.
func (v *LatencyVar) Set(latency time.Duration) { v.val.Store(int64(latency)) }

// Generate implements LatencyGenerator.
func (v *LatencyVar) Generate() time.Duration { return time.Duration(v.val.Load()) }

type JitterVar struct {
	val atomic.Value
}

func (v *JitterVar) Set(jitter time.Duration) { v.val.Store(jitter) }
func (v *JitterVar) Jitter() time.Duration    { return v.val.Load().(time.Duration) }

type LossVar struct {
	// See: https://github.com/golang/go/issues/21996
	val atomic.Uint64
}

func (v *LossVar) Set(loss float64) { v.val.Store(math.Float64bits(loss)) }
func (v *LossVar) Loss() float64    { return math.Float64frombits(v.val.Load()) }
