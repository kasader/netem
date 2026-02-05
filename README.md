# netem
A small Go network emulation package for "net" package interfaces. Allows for dynamic thread-safe configuration of connection bandwidth, latency, jitter, and packet loss via indirection.

This package was partially inspired by [cevatbarisyilmaz/lossy][1], and iterates upon its design aiming for:

1. A highly idiomatic and easy to use API for testing.
2. Adding the ability to alter the reading side of the connection (not just)
   the writing side.
3. Unbounded resource growth. At high throughput we spawn a goroutine per write/packet
   which could be better managed by a single queue. Even if goroutines are cheap they
   are not free.
4. Fixing major bugs in the implementation, notably on the stream-oriented net.Conn
   implementation side. There is stream corruption because Write() does not preserve
   the write order. In simulating a stream-oriented protocol a high jitter should
   cause [Head-of-Line Blocking][2], and NOT reorder the byte-stream.

    ```mermaid
    sequenceDiagram
    participant App as Application
    participant GA as Goroutine A (Data: "Hello")
    participant GB as Goroutine B (Data: "World")
    participant Sock as Underlying Socket

    Note over App: App sends "Hello" then "World"
    
    App->>GA: Write("Hello")
    App->>GB: Write("World")
    
    Note over GA: Processing... <br/>(Latency: 100ms)
    Note over GB: Processing... <br/>(Latency: 10ms)
    
    rect rgb(200, 50, 50, .1)
        Note right of GB: Jitter causes B to finish first
        GB->>Sock: Writes "World"
        Note over Sock: Buffer: "World"
        
        GA->>Sock: Writes "Hello"
        Note over Sock: Buffer: "WorldHello"
    end

    Note over Sock: !!! STREAM CORRUPTION: <br/>Expected "HelloWorld", Got "WorldHello"
    ```

This project uses golangci-lint. Run using `golangci-lint run ./...`. Maybe I will add a Makefile at some point (not sure).

[1]: https://github.com/cevatbarisyilmaz/lossy "cevatbarisyilmaz/lossy"
[2]: https://en.wikipedia.org/wiki/Head-of-line_blocking "Head-of-Line Blocking"
