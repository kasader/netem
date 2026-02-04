package netem

import "time"

// Option configures a netem connection.
type Option func(*config)

// WithLatency sets the base latency.
func WithLatency(d time.Duration) Option {
	return func(c *config) { c.latency = d }
}

// WithJitter sets the jitter (random variation) added to latency.
func WithJitter(d time.Duration) Option {
	return func(c *config) { c.jitter = d }
}

// WithPacketLoss sets the probability of dropping a packet (0.0 to 1.0).
func WithPacketLoss(rate float64) Option {
	return func(c *config) { c.lossRate = rate }
}

// WithBandwidth sets the bandwidth limit in bits per second.
func WithBandwidth(bps uint64) Option {
	return func(c *config) { c.bandwidth = bps }
}
