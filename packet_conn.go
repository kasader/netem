package netem

import (
	"net"
	"sync"
	"time"
)

// --- [net.PacketConn] implementation

// PacketConn TODO: insert doc.
type PacketConn struct {
	net.PacketConn

	p          PacketProfile
	headerSize int

	stopOnce sync.Once
	stopCh   chan struct{}
}

// NewPacketConn TODO: insert doc.
func NewPacketConn(c net.PacketConn, p PacketProfile) net.PacketConn {
	return &PacketConn{
		PacketConn: c,
		p:          p,
		headerSize: getHeaderSize(c.LocalAddr()),

		stopCh: make(chan struct{}),
	}
}

// Close implements net.PacketConn.
func (c *PacketConn) Close() error {
	c.stopOnce.Do(func() {
		if c.stopCh != nil {
			close(c.stopCh)
		}
	})
	return c.PacketConn.Close()
}

// LocalAddr implements net.PacketConn.
func (c *PacketConn) LocalAddr() net.Addr {
	panic("unimplemented")
}

// ReadFrom implements net.PacketConn.
func (*PacketConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	_ = p
	panic("unimplemented")
}

// SetDeadline implements net.PacketConn.
func (c *PacketConn) SetDeadline(t time.Time) error {
	_ = t
	panic("unimplemented")
}

// SetReadDeadline implements net.PacketConn.
func (c *PacketConn) SetReadDeadline(t time.Time) error {
	_ = t
	panic("unimplemented")
}

// SetWriteDeadline implements net.PacketConn.
func (c *PacketConn) SetWriteDeadline(t time.Time) error {
	_ = t
	panic("unimplemented")
}

// WriteTo implements net.PacketConn.
func (*PacketConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	_, _ = p, addr
	panic("unimplemented")
}

var _ net.PacketConn = (*PacketConn)(nil)
