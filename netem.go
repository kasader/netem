package netem

import (
	"net"
	"time"
)

const (
	// IPv4HeaderSize is the min size of an IPv4 header in bytes.
	IPv4HeaderSize = 20
	// IPv6HeaderSize is the fixed size of an IPv6 header in bytes.
	IPv6HeaderSize = 40
)

func getHeaderSize(addr net.Addr) int {
	var ip net.IP
	switch v := addr.(type) {
	case *net.UDPAddr:
		ip = v.IP
	case *net.TCPAddr:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	// Best-effort to parse custom implementation (if provided).
	default:
		host, _, err := net.SplitHostPort(addr.String())
		if err != nil {
			// I guess just try parsing directly if SplitHostPort fails...?
			ip = net.ParseIP(addr.String())
		} else {
			ip = net.ParseIP(host)
		}
	}
	// Determine our header overhead from the [net.IP].
	overhead := IPv6HeaderSize // Assume worst case (IPv6)
	if ip != nil && ip.To4() != nil {
		overhead = IPv4HeaderSize
	}
	return overhead
}

func transmissionTime(bandwidth uint64, size, overhead int) time.Duration {
	if bandwidth == 0 {
		return 0
	}
	// Convert byte-size to bit-size (we measure in bits/second).
	totalBits := float64(size+overhead) * 8.0

	seconds := totalBits / float64(bandwidth)
	return time.Duration(seconds * float64(time.Second))
}

func delayTime(latency Latency, jitter Jitter) time.Duration {
	var totalDelay time.Duration
	if latency != nil {
		totalDelay += latency.Duration()
	}
	if jitter != nil {
		totalDelay += jitter.Duration()
	}
	return totalDelay
}
