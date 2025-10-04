# HLA=tcp

> Hidden Lake Adapter (TCP)

<img src="images/hla_tcp_logo.png" alt="hla_tcp_logo.png"/>

The `Hidden Lake Adapter (TCP)` allows adapt HL traffic based on the TCP protocol.

## Installation

```bash
$ go install github.com/number571/hidden-lake/cmd/hla/hla-tcp@latest
```

## How it works

HLA=tcp uses `go-peer` (pkg/network) implementation.

## Supported platforms

- Windows (x86_64, arm64)
- Linux (x86_64, arm64)
- MacOS (x86_64, arm64)

## Build and run

Default build and run

```bash 
$ go run ./cmd/hla/hla-tcp

> [INFO] 2023/06/03 15:30:31 HLA=tcp is running...
> ...
```

Open port `9521` (TCP, external), `9522` (HTTP, internal).
Creates [`./hla-tcp.yml`](./hla-tcp.yml) file.

## Running options

```bash
$ hla-tcp --path /root
# path = path to config
```
