package netem

// LinkProfile defines the shared physical properties of a network link.
type LinkProfile struct {
	Latency   Latency
	Jitter    Jitter
	Bandwidth Bandwidth
}

// PacketProfile extends the link with datagram-specific behaviors.
type PacketProfile struct {
	LinkProfile
	Loss Loss
}

// StreamProfile extends the link with stream-specific behaviors.
type StreamProfile struct {
	LinkProfile
	Fault Fault
}
