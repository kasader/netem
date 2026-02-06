# netem

[![Go Reference](https://pkg.go.dev/badge/github.com/kasader/netem/netem.svg)](https://pkg.go.dev/github.com/kasader/netem)
[![Go Report Card](https://goreportcard.com/badge/github.com/kasader/netem/netem)](https://goreportcard.com/report/github.com/kasader/netem)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![GitHub Release](https://img.shields.io/github/v/release/kasader/netem?include_prereleases)](https://github.com/kasader/releases)

`netem` is a lightweight Go network emulation package designed to wrap standard `net.Conn` and `net.PacketConn` interfaces. It allows developers to simulate real-world network conditions, like restricted bandwidth, high latency, jitter, and packet loss, directly within their Go tests.

## Why netem?

While several tools exist for network emulation at the OS level (like `tc-netem` on Linux), `netem` brings this capability directly into the Go test suite. It allows you to wrap existing connections to simulate poor network conditions without requiring administrative privileges or external dependencies.

### Key Improvements

This package iterates on the design of [`cevatbarisyilmaz/lossy`][1], focusing on:

- **Idiomatic API**: Designed to feel like a natural extension of the `net` package.

- **Resource Efficiency**: Replaces "goroutine-per-packet" models with internal FIFO queues to prevent unbounded resource growth at high throughput.

- **Stream Integrity**: Specifically addresses the "WorldHello" corruption bug. In stream-oriented protocols (TCP), high jitter should cause [Head-of-Line Blocking][2], not out-of-order byte delivery.

## Dynamic Configuration (Inspired by `slog`)

One of the standout features of `netem` is its thread-safe indirection system, inspired by the design of `slog.LevelVar`.

Instead of static values, `netem` uses **Policies**. By using types like `policy.LatencyVar` or `policy.BandwidthVar`, you can alter network conditions on the fly for an *active* connection without needing to reconnect or use complex synchronization in your application code.

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

## Core Features (v0.1.0)

### 1. Dynamic Policies (`netem/policy`)

Inspired by the design of `slog.LevelVar`, `netem` uses an indirection system for its metrics. By using `Var` types, you can change network conditions on-the-fly for active connections.

- **Bandwidth**: Limit throughput (bits per second) via `StaticBandwidth` or the thread-safe `BandwidthVar`.

- **Latency**: Add base propagation delay.

- **Jitter**: Introduce variance in delivery time. `RandomJitter` provides amplitude-based variance.

- **Loss**: Simulate packet drops with `RandomLoss`.

- **Fault**: Manually trigger connection failures or closures.

### 2. Protocol-Specific Wrappers

- **`Conn` (Stream-Oriented)**: Designed for TCP-like behavior. Uses an internal FIFO queue to ensure that jitter manifests as **Head-of-Line Blocking** rather than byte-stream corruption.

- **`PacketConn` (Packet-Oriented)**: Designed for UDP-like behavior. Supports natural packet reordering, allowing later datagrams to overtake earlier ones if jitter is sufficiently high.


[1]: https://github.com/cevatbarisyilmaz/lossy "cevatbarisyilmaz/lossy"
[2]: https://en.wikipedia.org/wiki/Head-of-line_blocking "Head-of-Line Blocking"
