// Package netem provides network emulation wrappers for the [net.PacketConn]
// and [net.Conn] interfaces; it allows for dynamic thread-safe configuration
// of connection bandwidth, latency, jitter, and packet loss via indirection.
package netem

import (
	"net"
	"sync"
	"time"
)

// ConfigurationVar is a [config] variable, to allow an emulated
// connection to change dynamically.
//
// It implements [...] as well as Set methods, and it is safe for use by
// multiple goroutines. The zero value corresponds to [...].
type ConfigurationVar struct {
	mu  sync.RWMutex
	cfg *Config
}

type Config struct {
	Jitter  JitterGenerator
	Latency LatencyGenerator

	Bandwidth uint64 // bits per second, 0 = infinite
	LossRate  float64
}

func defaultConfig() *Config {
	return &Config{
		LossRate:  0,
		Bandwidth: 0,
	}
}

// --- [net.PacketConn] implementation

func NewPacketConn(c net.PacketConn, opts ...Option) net.PacketConn {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	return &PacketConn{
		PacketConn: c,
		config:     cfg,
	}
}

type PacketConn struct {
	net.PacketConn
	config *Config
}

// Close implements net.PacketConn.
func (p *PacketConn) Close() error {
	panic("unimplemented")
}

// LocalAddr implements net.PacketConn.
func (p *PacketConn) LocalAddr() net.Addr {
	panic("unimplemented")
}

// ReadFrom implements net.PacketConn.
func (*PacketConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	panic("unimplemented")
}

// SetDeadline implements net.PacketConn.
func (p *PacketConn) SetDeadline(t time.Time) error {
	panic("unimplemented")
}

// SetReadDeadline implements net.PacketConn.
func (p *PacketConn) SetReadDeadline(t time.Time) error {
	panic("unimplemented")
}

// SetWriteDeadline implements net.PacketConn.
func (p *PacketConn) SetWriteDeadline(t time.Time) error {
	panic("unimplemented")
}

// WriteTo implements net.PacketConn.
func (*PacketConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	panic("unimplemented")
}

var _ net.PacketConn = (*PacketConn)(nil)

// --- [net.Conn] implementation

type Conn struct {
	net.Conn
	// ...
}

// Close implements net.Conn.
func (c *Conn) Close() error {
	panic("unimplemented")
}

// LocalAddr implements net.Conn.
func (c *Conn) LocalAddr() net.Addr {
	panic("unimplemented")
}

// Read implements net.Conn.
func (c *Conn) Read(b []byte) (n int, err error) {
	panic("unimplemented")
}

// RemoteAddr implements net.Conn.
func (c *Conn) RemoteAddr() net.Addr {
	panic("unimplemented")
}

// SetDeadline implements net.Conn.
func (c *Conn) SetDeadline(t time.Time) error {
	panic("unimplemented")
}

// SetReadDeadline implements net.Conn.
func (c *Conn) SetReadDeadline(t time.Time) error {
	panic("unimplemented")
}

// SetWriteDeadline implements net.Conn.
func (c *Conn) SetWriteDeadline(t time.Time) error {
	panic("unimplemented")
}

// Write implements net.Conn.
func (c *Conn) Write(b []byte) (n int, err error) {
	panic("unimplemented")
}

var _ net.Conn = (*Conn)(nil)
