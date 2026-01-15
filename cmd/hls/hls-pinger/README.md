# HLS=pinger

> Hidden Lake Service (Pinger)

<img src="images/hls_pinger_logo.png" alt="hls_pinger_logo.png"/>

The `Hidden Lake Service (Pinger)` allows you to ping a node to see if a participant is online.

## Installation

```bash
$ go install github.com/number571/hidden-lake/cmd/hls/hls-pinger@latest
```

## How it works

HLS=pinger is the simplest service

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
$ go run ./cmd/hls/hls-pinger

> [INFO] 2023/06/03 15:30:31 HLS=pinger is running...
> ...
```

Open ports `9551` (HTTP, internal), `9552` (HTTP, incoming).
Creates [`./hls-pinger.yml`](./hls-pinger.yml) file.

## Running options

```bash
$ hls-pinger --path /root
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
$ ./_request/raw/request.sh
```

Got response
```json
{"code":200,"head":{"Content-Type":"text/plain"}}
```

## HLS API

```
1. GET /api/index
2. GET /api/send/ping
```

### 1. /api/index

#### 1.1. GET Request

```bash
curl -i -X GET http://localhost:9551/api/index
```

#### 1.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Thu, 15 Jan 2026 10:30:39 GMT
Content-Length: 18

hidden-lake-kernel
```

### 2. /api/send/ping

#### 2.1. GET Request

```bash
curl -i -X GET "http://localhost:9551/api/send/ping?friend=Bob"
```

#### 2.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Thu, 15 Jan 2026 10:47:04 GMT
Content-Length: 0
```
