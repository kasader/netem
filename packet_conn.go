package netem

import (
	"net"
	"time"
)

// --- [net.PacketConn] implementation

func NewPacketConn(c net.PacketConn, p PacketProfile) net.PacketConn {
	panic("unimplemented")
}

type PacketConn struct {
	net.PacketConn
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
