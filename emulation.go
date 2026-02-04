package netem

import (
	"math"
	"sync/atomic"
	"time"
)

type BandwidthProvider interface{ Bandwidth() uint64 }

type LatencyProvider interface{ Latency() time.Duration }

type JitterProvider interface{ Jitter() time.Duration }

type LossProvider interface{ Loss() float64 }

type BandwidthVar struct {
	val atomic.Uint64
}

func (v *BandwidthVar) Set(bandwidth uint64) { v.val.Store(bandwidth) }
func (v *BandwidthVar) Bandwidth() uint64    { return v.val.Load() }

type LatencyVar struct {
	val atomic.Value
}

func (v *LatencyVar) Set(latency time.Duration) { v.val.Store(latency) }
func (v *LatencyVar) Latency() time.Duration    { return v.val.Load().(time.Duration) }

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
