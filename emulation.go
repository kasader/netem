package netem

import (
	"math"
	"sync/atomic"
	"time"
)

// LatencyGenerator calculates the bandwidth for a packet.
type BandwidthGenerator interface{ Generate() uint64 }

// LatencyGenerator calculates the latency for a packet.
type LatencyGenerator interface{ Generate() time.Duration }

// JitterGenerator calculates the jitter for a packet.
type JitterGenerator interface{ Generate() time.Duration }

// LatencyGenerator calculates the loss for a packet.
type LossGenerator interface{ Generate() float64 }

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
