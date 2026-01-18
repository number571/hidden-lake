# HLS=messenger

> Hidden Lake Service (Messenger)

<img src="images/hls_messenger_logo.png" alt="hls_messenger_logo.png"/>

The `Hidden Lake Service (Messenger)` is a messenger based on the core of an anonymous network with theoretically provable anonymity of HLK. A feature of this messenger is the provision of anonymity of the fact of transactions (sending, receiving).

> More information about HLS=messenger in the [habr.com/ru/post/701488](https://habr.com/ru/post/701488/ "Habr HLS=messenger")

## Installation

```bash
$ go install github.com/number571/hidden-lake/cmd/hls/hls-messenger@latest
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
$ go run ./cmd/hls/hls-messenger

> [INFO] 2023/06/03 15:30:31 HLS=messenger is running...
> ...
```

Open ports `9591` (HTTP, internal) and `9592` (HTTP, incoming).
Creates [`./hls-messenger.yml`](./hls-messenger.yml) and `./hls-messenger.db` files.

## Running options

```bash
$ hls-messenger --path /root
# path = path to config and database files
```

## Example

The example will involve (as well as in HLK) five nodes `node1_hlm, node2_hlm` and `middle_hla_tcp_1, middle_hla_tcp_2, middle_hla_tcp_3`. The three `HLA=tcp` nodes are only needed for communication between `node1_hlm` and `node2_hlm` nodes. Each of the remaining ones is a combination of HLK and HLS=messenger, where HLS=messenger plays the role of an application and services (as it was depicted in `Figure 3` HLK readme).

Build and run nodes
```bash
$ cd examples/messenger/routing
$ make
```

Than run commands
```bash
### Terminal 1 ###
$ cd examples/messenger
$ make request-node1
### Terminal 2 ###
$ cd examples/messenger
$ make request-node2
```

## HLS API

```
1. GET  /api/index              | params = [] 
                                |> description = get name of service
2. POST /api/chat/message       | params = ["friend":string]
                                |> description = send message to chat
3. GET  /api/chat/history/load  | params = ["friend":string,"start":uint64,"count":uint64,"select":string]
                                |> description = get list of messages from chat
4. GET  /api/chat/history/size  | params = ["friend":string]
                                |> description = get count of messages in the chat
5. GET  /api/chat/subscribe     | params = ["friend":string]
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

hidden-lake-service=messenger
```

### 2. /api/chat/message

#### 2.1. POST Request

```bash
curl -i -X POST "http://localhost:9591/api/chat/message?friend=Bob" --data 'hello, world!'
```

#### 2.1. POST Response

```
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Fri, 16 Jan 2026 18:55:12 GMT
Content-Length: 19

2026-01-16T18:55:12
```

### 3. /api/chat/history/load

#### 3.1. GET Request

```bash
curl -i -X GET "http://localhost:9591/api/chat/history/load?friend=Bob&start=0&count=10&select=asc"
```

#### 3.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 16 Jan 2026 18:56:08 GMT
Content-Length: 69

[{"incoming":false,"message":"hello, world!","timestamp":1768589712}]
```

### 4. /api/chat/history/size

#### 4.1. GET Request

```bash
curl -i -X GET "http://localhost:9591/api/chat/history/size?friend=Bob"
```

#### 4.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 16 Jan 2026 18:56:08 GMT
Content-Length: 1

1
```

### 5. /api/chat/subscribe

#### 5.1. GET Request

```bash
curl -i -X GET "http://localhost:9591/api/chat/subscribe?friend=Bob"
```

#### 5.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Fri, 16 Jan 2026 18:57:22 GMT
Content-Length: 58

{"incoming":true,"message":"hello","timestamp":1768589838}
```
