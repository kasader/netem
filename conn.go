package netem

import (
	"net"
	"sync"
	"time"
)

// --- [net.Conn] implementation

// Conn TODO: insert docs.
type Conn struct {
	net.Conn

	p          StreamProfile
	headerSize int

	mu         sync.Mutex
	throttleMu sync.Mutex

	stopOnce sync.Once
	stopCh   chan struct{}
}

// NewConn TODO: insert docs.
func NewConn(c net.Conn, p StreamProfile) net.Conn {
	return &Conn{
		Conn:       c,
		p:          p,
		headerSize: getHeaderSize(c.LocalAddr()),

		stopCh: make(chan struct{}),
	}
}

// Close implements net.Conn.
func (c *Conn) Close() error {
	c.stopOnce.Do(func() {
		if c.stopCh != nil {
			close(c.stopCh)
		}
	})
	return c.Conn.Close()
}

// LocalAddr implements net.Conn.
func (c *Conn) LocalAddr() net.Addr {
	panic("unimplemented")
}

// Read implements net.Conn.
func (c *Conn) Read(b []byte) (n int, err error) {
	_ = b
	panic("unimplemented")
}

// RemoteAddr implements net.Conn.
func (c *Conn) RemoteAddr() net.Addr {
	panic("unimplemented")
}

// SetDeadline implements net.Conn.
func (c *Conn) SetDeadline(t time.Time) error {
	_ = t
	panic("unimplemented")
}

// SetReadDeadline implements net.Conn.
func (c *Conn) SetReadDeadline(t time.Time) error {
	_ = t
	panic("unimplemented")
}

// SetWriteDeadline implements net.Conn.
func (c *Conn) SetWriteDeadline(t time.Time) error {
	_ = t
	panic("unimplemented")
}

// Write implements net.Conn.
func (c *Conn) Write(b []byte) (n int, err error) {
	_ = b
	select {
	case <-c.stopCh:
		return c.Conn.Write(b)
	default:
		// don't block on non-closed stopCh
	}

	// Since we return immediately, we must copy the data to own it.
	// Otherwise, the caller might modify 'b' while our goroutine is sleeping.
	data := make([]byte, len(b))
	copy(data, b)

	go func() {
		// Serialization Delay (Bandwidth)
		if c.p.Bandwidth != nil {
			limit := c.p.Bandwidth.Limit()
			if limit != 0 {
				delay := transmissionTime(limit, len(data), c.headerSize)
				c.throttleMu.Lock()
				time.Sleep(delay)
				c.throttleMu.Unlock()
			}
		}
		// Fault Injection (Fault)
		if c.p.Fault != nil && c.p.Fault.ShouldClose() {
			c.Close() // simulate abrupt connection drop
			return
		}
		// Propagation delay (Latency + Jitter)
		if d := delayTime(c.p.Latency, c.p.Jitter); d != 0 {
			time.Sleep(d)
		}
		_, _ = c.Conn.Write(b)
	}()
	return len(b), nil
}

var _ net.Conn = (*Conn)(nil)
