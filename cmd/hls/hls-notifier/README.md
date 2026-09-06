# HLS=notifier

> Hidden Lake Service (Notifier)

<img src="images/hls_notifier_logo.png" alt="hls_notifier_logo.png"/>

The `Hidden Lake Service (Notifier)` is a notifier based on the core of an anonymous network with theoretically provable anonymity of HLK. A feature of this notifier is the provision of anonymity of the fact of transactions (sending, receiving) text messages.

> More information about HLS=notifier in the [habr.com/ru/post/701488](https://habr.com/ru/post/701488/ "Habr HLS=notifier")

## Installation

```bash
$ go install github.com/number571/hidden-lake/cmd/hls/hls-notifier@latest
```

## How it works

Most of the code is a call to API functions from the HLK kernel. Thanks to this approach, implicit authorization of users is formed from the state of the anonymizing service.

## Supported platforms

- Windows (x86_64, arm64)
- Linux (x86_64, arm64)
- MacOS (x86_64, arm64)

## Build and run

Default build and run

```bash 
$ go run ./cmd/hls/hls-notifier

> [INFO] 2023/06/03 15:30:31 HLS=notifier is running...
> ...
```

Open ports `9591` (HTTP, internal) and `9592` (HTTP, incoming).
Creates [`./hls-notifier.yml`](./hls-notifier.yml) and `./hls-notifier.db` files.

## Running options

```bash
$ hls-notifier --path /root
# path = path to config and database files
```

## Example

The example will involve (as well as in HLK) five nodes `node1_hlm, node2_hlm` and `middle_hla_tcp_1, middle_hla_tcp_2, middle_hla_tcp_3`. The three `HLA=tcp` nodes are only needed for communication between `node1_hlm` and `node2_hlm` nodes. Each of the remaining ones is a combination of HLK and HLS=notifier, where HLS=notifier plays the role of an application and services (as it was depicted in `Figure 3` HLK readme).

Build and run nodes
```bash
$ cd examples/notifier/routing
$ make
```

Than run commands
```bash
### Terminal 1 ###
$ cd examples/notifier
$ make request-node1
### Terminal 2 ###
$ cd examples/notifier
$ make request-node2
```

## HLS API

```
1. GET      /api/index          | params = [] 
                                |> description = get name of service
2. GET/POST /api/chat/push      | params = ["friend":string]
                                |> description = get limit message size / send message to chat
3. GET      /api/chat/load      | params = ["friend":string,"index":uint64]
                                |> description = get message from chat by index
4. GET      /api/chat/size      | params = ["friend":string]
                                |> description = get count of messages in the chat
5. GET      /api/chat/listen    | params = ["friend":string]
                                |> description = try get message from chat with longpoll method 
```

### 1. /api/index

#### 1.1. GET Request

```bash
curl -i -X GET http://localhost:9591/api/index
```

#### 1.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Fri, 16 Jan 2026 18:54:07 GMT
Content-Length: 29

hidden-lake-service=notifier
```

### 2. /api/chat/push

#### 2.1. GET Request

```bash
curl -i -X GET "http://localhost:9591/api/chat/push"
```

#### 2.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Sun, 18 Jan 2026 20:38:58 GMT
Content-Length: 4

3552
```

#### 2.2. POST Request

```bash
curl -i -X POST "http://localhost:9591/api/chat/push?friend=Bob" --data 'hello, world!'
```

#### 2.2. POST Response

```
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Fri, 16 Jan 2026 18:55:12 GMT
Content-Length: 19

2026-01-16T18:55:12
```

### 3. /api/chat/load

#### 3.1. GET Request

```bash
curl -i -X GET "http://localhost:9591/api/chat/load?friend=Bob&index=0"
```

#### 3.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 16 Jan 2026 18:56:08 GMT
Content-Length: 69

{"incoming":false,"message":"hello, world!","timestamp":1768589712}
```

### 4. /api/chat/size

#### 4.1. GET Request

```bash
curl -i -X GET "http://localhost:9591/api/chat/size?friend=Bob"
```

#### 4.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 16 Jan 2026 18:56:08 GMT
Content-Length: 1

1
```

### 5. /api/chat/listen

#### 5.1. GET Request

```bash
curl -i -X GET "http://localhost:9591/api/chat/listen?friend=Bob"
```

#### 5.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Fri, 16 Jan 2026 18:57:22 GMT
Content-Length: 58

{"incoming":true,"message":"hello","timestamp":1768589838}
```
