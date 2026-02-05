package netem

import (
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type writeReq struct {
	data []byte
	due  time.Time
}

type Conn struct {
	net.Conn
	p          StreamProfile
	headerSize int

	writeCh chan writeReq // FIFO queue

	mu         sync.Mutex
	throttleMu sync.Mutex

	stopOnce sync.Once
	stopCh   chan struct{}

	writeDeadline atomic.Value
}

// NewConn TODO: insert docs.
func NewConn(c net.Conn, p StreamProfile) net.Conn {
	nc := &Conn{
		Conn: c,
		p:    p,

		// Buffered to allow bursting.
		// TODO: Should the WriteCh length be configurable?
		writeCh: make(chan writeReq, 1024),
		stopCh:  make(chan struct{}),

		headerSize: getHeaderSize(c.LocalAddr()),
	}
	nc.writeDeadline.Store(time.Time{})
	go nc.linkLoop()
	return nc
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

// SetDeadline implements net.Conn.
func (c *Conn) SetDeadline(t time.Time) error {
	c.writeDeadline.Store(t)
	return c.Conn.SetDeadline(t)
}

// SetWriteDeadline implements net.Conn.
func (c *Conn) SetWriteDeadline(t time.Time) error {
	c.writeDeadline.Store(t)
	return c.Conn.SetWriteDeadline(t)
}

// Write implements net.Conn.
func (c *Conn) Write(b []byte) (n int, err error) {
	if c.isWriteDeadline() {
		return 0, os.ErrDeadlineExceeded
	}
	// 1. Calculate packet delays immediately.
	serializationDelay := transmissionTime(c.p.Bandwidth, len(b), c.headerSize)
	propagationDelay := delayTime(c.p.Latency, c.p.Jitter)

	// 2. Schedule our packet using the calculated delay.
	req := writeReq{
		data: make([]byte, len(b)),
		due:  time.Now().Add(serializationDelay + propagationDelay),
	}
	copy(req.data, b)

	// 3. Enqueue (and do a non-blocking check for stop).
	select {
	case <-c.stopCh:
		return c.Conn.Write(b)
	case c.writeCh <- req:
		return len(b), nil
	}
}

var _ net.Conn = (*Conn)(nil)

// Handles writes in strict order.
func (c *Conn) linkLoop() {
	for {
		select {
		case <-c.stopCh:
			return
		case req := <-c.writeCh:
			// 1. Perform fault injection before writing.
			if c.p.Fault != nil && c.p.Fault.ShouldClose() {
				c.Close()
			}
			// 2. Wait until due time.
			wait := time.Until(req.due)
			if wait > 0 {
				time.Sleep(wait)
			}
			// 3. Write; and because we pull from the channel we can
			// assume that our packets must be written in order.
			c.Conn.Write(req.data)
		}
	}
}

func (c *Conn) isWriteDeadline() bool {
	wdl := c.writeDeadline.Load().(time.Time)
	return !wdl.IsZero() && wdl.Before(time.Now())
}
