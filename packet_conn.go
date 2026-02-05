package netem

import (
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// packetReq holds the data and the scheduled arrival time.
type packetReq struct {
	data []byte
	addr net.Addr
	due  time.Time
}

// packetHeap is a Min-Heap sorted by 'due' time.
type packetHeap []packetReq

func (h packetHeap) Len() int           { return len(h) }
func (h packetHeap) Less(i, j int) bool { return h[i].due.Before(h[j].due) }
func (h packetHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *packetHeap) Push(x any) {
	*h = append(*h, x.(packetReq))
}

func (h *packetHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// PacketProfile extends the link with datagram-specific behaviors.
type PacketProfile struct {
	// MTU (Maximum Transmission Unit) is the largest packet size allowed.
	// This value includes L3/L4 headers.
	//
	// Defaults to [EthernetDefaultMTU] if 0.
	MTU uint

	Latency   Latency
	Jitter    Jitter
	Bandwidth Bandwidth
	Loss      Loss
}

// PacketConn TODO: insert doc.
type PacketConn struct {
	net.PacketConn
	headerSize int
	mss        int
	p          PacketProfile

	writeDeadline atomic.Value

	stopOnce sync.Once
	stopCh   chan struct{}
}

// NewPacketConn TODO: insert doc.
func NewPacketConn(c net.PacketConn, p PacketProfile) net.PacketConn {
	headerSize := getHeaderSize(c.LocalAddr())
	mtu := p.MTU
	if mtu == 0 {
		mtu = EthernetDefaultMTU
	}
	// Enforce minimum mss.
	mss := max(1, int(mtu)-headerSize)

	nc := &PacketConn{
		PacketConn: c,
		headerSize: getHeaderSize(c.LocalAddr()),
		mss:        mss,
		p:          p,

		stopCh: make(chan struct{}),
	}
	nc.writeDeadline.Store(time.Time{})
	return nc
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

// SetDeadline implements net.PacketConn.
func (c *PacketConn) SetDeadline(t time.Time) error {
	c.writeDeadline.Store(t)
	return c.PacketConn.SetDeadline(t)
}

// SetWriteDeadline implements net.PacketConn.
func (c *PacketConn) SetWriteDeadline(t time.Time) error {
	c.writeDeadline.Store(t)
	return c.PacketConn.SetWriteDeadline(t)
}

// WriteTo implements net.PacketConn.
func (*PacketConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	_, _ = p, addr
	panic("unimplemented")
}

var _ net.PacketConn = (*PacketConn)(nil)

func (c *PacketConn) isWriteDeadline() bool {
	wdl := c.writeDeadline.Load().(time.Time)
	return !wdl.IsZero() && wdl.Before(time.Now())
}
