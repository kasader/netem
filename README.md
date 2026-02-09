# netem

[![Go Reference](https://pkg.go.dev/badge/github.com/kasader/netem/netem.svg)](https://pkg.go.dev/github.com/kasader/netem)
[![Go Report Card](https://goreportcard.com/badge/github.com/kasader/netem)](https://goreportcard.com/report/github.com/kasader/netem)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![GitHub Release](https://img.shields.io/github/v/release/kasader/netem?include_prereleases)](https://github.com/kasader/netem/releases)

`netem` is a lightweight Go network emulation package designed to wrap standard `net.Conn` and `net.PacketConn` interfaces. It allows simulation of real-world network conditions, like restricted bandwidth, high latency, jitter, and packet loss, directly within Go tests.

## Why netem?

While several tools exist for network emulation at the OS level (like `tc-netem` on Linux), `netem` provides this capability to Go testing. It allows you to wrap existing connections to simulate poor network conditions without requiring non-Go dependencies.

### Key Improvements

This package iterates on the design of [`cevatbarisyilmaz/lossy`][1], focusing on:

- **Idiomatic API**: Designed to feel like a natural extension of the `net` package.

- **Resource Efficiency**: Replaces the "goroutine-per-packet" models to prevent unbounded resource growth at high throughput.

- **Stream Integrity**: Addresses the "World!Hello, " corruption bug. I.e, in stream-oriented protocols (TCP), high jitter should cause [Head-of-Line Blocking][2], and not out-of-order byte delivery (as the other package does).

## Dynamic Configuration (Inspired by `slog`)

Thread-safe indirection is usable out-of-the-box, inspired by the design of `slog.LevelVar`.

Instead of static values, `netem` uses **Policies**. By using types like `policy.LatencyVar` or `policy.BandwidthVar`, you can alter network conditions on the fly for an *active* connection.

## Usage

```go
import (
    "github.com/kasader/netem"
    "github.com/kasader/netem/policy"
)

func TestMyNetworkCode(t *testing.T) {
    // 1. Define a profile
    lat := &policy.LatencyVar{}
    lat.Set(100 * time.Millisecond)

    profile := netem.PacketProfile{
        Latency: lat,
        Jitter:  policy.RandomJitter(20 * time.Millisecond),
        Loss:    policy.RandomLoss(0.01), // 1% loss
    }

    // 2. Wrap an existing connection
    conn := netem.NewPacketConn(rawUDPConn, profile)

    // 3. Dynamically change conditions later
    lat.Set(500 * time.Millisecond) 
}
```

[1]: https://github.com/cevatbarisyilmaz/lossy "cevatbarisyilmaz/lossy"
[2]: https://en.wikipedia.org/wiki/Head-of-line_blocking "Head-of-Line Blocking"
