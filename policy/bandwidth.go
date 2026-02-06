package policy

import (
	"sync/atomic"
)

// BandwidthFunc enables a simple function to satisfy the [Bandwidth] interface.
type BandwidthFunc func() uint64

// Limit implements the [Bandwidth] interface.
func (f BandwidthFunc) Limit() uint64 { return f() }

// StaticBandwidth returns a constant throughput.
func StaticBandwidth(bps uint64) BandwidthFunc {
	return BandwidthFunc(func() uint64 { return bps })
}

// BandwidthVar is a thread-safe, mutable [Bandwidth] provider.
// It allows you to change the bandwidth of a running simulation.
type BandwidthVar struct{ val atomic.Uint64 }

// Set updates the bandwidth safely.
func (v *BandwidthVar) Set(bandwidth uint64) { v.val.Store(bandwidth) }

// Limit implements the [Bandwidth] interface.
func (v *BandwidthVar) Limit() uint64 { return v.val.Load() }
