# HLR

> Hidden Lake Pinger

<img src="images/hlp_logo.png" alt="hlp_logo.png"/>

The `Hidden Lake Pinger` allows you to ping a node to see if a participant is online.

## Installation

```bash
$ go install github.com/number571/hidden-lake/cmd/hlp@latest
```

## How it works

HLP is the simplest service

```go
func handler() {
    if method != GET {
        response(StatusMethodNotAllowed)
        return
    }
    response(StatusOK)
}
```

## Supported platforms

- Windows (x86_64, arm64)
- Linux (x86_64, arm64)
- MacOS (x86_64, arm64)

## Build and run

Default build and run

```bash 
$ go run ./cmd/hlp

> [INFO] 2023/06/03 15:30:31 HLP is running...
> ...
```

Open port `9552` (HTTP, incoming).
Creates [`./hlp.yml`](./hlp.yml) file.

## Running options

```bash
$ hlp --path /root
# path = path to config
```

## Example

The example will involve two nodes `recv_hlc, send_hls` and three repeaters `middle_hla_tcp_1, middle_hla_tcp_2, middle_hla_tcp_3`. The three remaining nodes are used only for the successful connection of the two main nodes. In other words, `HLA=tcp` nodes are traffic relay nodes.

Build and run nodes
```bash
$ cd examples/pinger/routing
$ make
```

Than run command
```bash
$ cd examples/pinger
$ make request # go run ./_request/main.go
```

Got response
```json
{"code":200,"head":{"Content-Type":"text/plain"}}
```
