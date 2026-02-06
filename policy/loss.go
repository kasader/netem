package policy

import (
	"math"
	"math/rand/v2"
	"sync/atomic"
)

// LossFunc enables a simple function to satisfy the [Loss] interface.
type LossFunc func() bool

// Drop implements the [Loss] interface.
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

// Drop implements the [Loss] interface.
func (v *LossVar) Drop() bool {
	rate := math.Float64frombits(v.val.Load())
	return RandomLoss(rate)()
}
