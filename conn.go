package netem

import (
	"net"
	"sync"
	"time"
)

// --- [net.Conn] implementation

type Conn struct {
	net.Conn
	p StreamProfile

	stopOnce sync.Once
	stopCh   chan struct{}
}

func NewConn(c net.Conn, p StreamProfile) net.Conn {
	return &Conn{
		Conn: c,
		p:    p,

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
