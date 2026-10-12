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

## HLS/HLK API

### 1. /notify

- `200` - return "ok"
- `400` - invalid input params (body)
- `405` - invalid method (not POST)
- `409` - failed read HTTP request body
- `502` - get aliasName HLK

```json
// header request size = (4 + 56) = 60B
// body request size can be = (144 - 60) = 84B
[uint32]{"method":"POST","host":"hls-notifier","path":"/notify"}[body]

// header response size = (4 + 12) = 16B
// body response size can be = (144 - 16) = 128B
[uint32]{"code":200}[body]
```

### 2. /push

- `200` - return "ok"
- `202` - message already exist in database
- `400` - invalid input params (uid, body)
- `405` - invalid method (not POST)
- `409` - failed read HTTP request body
- `500` - load/push from database
- `502` - get aliasName HLK

```json
// header request size = (4 + 81) = 85B
// body request size can be = (144 - 85) = 59B
[uint32]{"method":"POST","host":"hls-notifier","path":"/push?uid=0000000000000000000000"}[body]

// header response size = (4 + 12) = 16B
// body response size can be = (144 - 16) = 128B
[uint32]{"code":200}[body]
```

## HLS API

```
1. GET          /api/index          | params = [] 
                                    |> description = get name of service
2. GET/POST     /api/chan/notify    | params = ["friend":string]
                                    |> description = get limit notify size / send notify to chan
3. GET/POST     /api/chan/push      | params = ["friend":string]
                                    |> description = get limit message size / send message to chan
4. GET          /api/chan/load      | params = ["friend":string,"index":uint64]
                                    |> description = get message from chan by index
5. GET          /api/chan/size      | params = ["friend":string]
                                    |> description = get count of messages in the chan
6. GET          /api/chan/listen    | params = ["friend":string]
                                    |> description = try get message from chan with longpoll method 
```

### 1. /api/index

- `200` - return "hidden-lake-service=notifier"

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

### 2. /api/chan/notify

- `200` - return limit size (GET) or timestamp (POST)
- `400` - invalid input params (friend, body)
- `405` - invalid method (not GET, POST)
- `409` - failed read HTTP request body
- `502` - failed send request (HLK side)

#### 2.1. GET Request

```bash
curl -i -X GET "http://localhost:9591/api/chan/notify"
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
curl -i -X POST "http://localhost:9591/api/chan/notify?friend=Bob&uid=0000000000000000000000" --data 'hello, world!'
```

#### 2.2. POST Response

```
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Fri, 16 Jan 2026 18:55:12 GMT
Content-Length: 19

2026-01-16T18:55:12
```

### 3. /api/chan/push

- `200` - return limit size (GET) or timestamp (POST)
- `202` - message already registered and delivered
- `208` - message in the fetch process
- `400` - invalid input params (friend, uid, body)
- `405` - invalid method (not GET, POST)
- `409` - failed read HTTP request body
- `418` - got bad response from node
- `500` - load/push/recv message into database
- `502` - failed send request (HLK side)
- `503` - message was sent but no response was received

#### 3.1. GET Request

```bash
curl -i -X GET "http://localhost:9591/api/chan/push"
```

#### 3.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Sun, 18 Jan 2026 20:38:58 GMT
Content-Length: 4

3525
```

#### 3.2. POST Request

```bash
curl -i -X POST "http://localhost:9591/api/chan/push?friend=Bob&uid=0000000000000000000000" --data 'hello, world!'
```

#### 3.2. POST Response

```
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Fri, 16 Jan 2026 18:55:12 GMT
Content-Length: 19

2026-01-16T18:55:12
```

### 4. /api/chan/load

- `200` - return message from database
- `400` - invalid input params (friend, index)
- `404` - message not found in database
- `405` - invalid method (not GET)
- `500` - load message from database

#### 4.1. GET Request

```bash
curl -i -X GET "http://localhost:9591/api/chan/load?friend=Bob&index=0"
```

#### 4.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 16 Jan 2026 18:56:08 GMT
Content-Length: 89

{"head":{"s":2,"u":"0000000000000000000000"},"body":{"d":"hello, world!","t":1768589712}}
```

### 5. /api/chan/size

- `200` - return count of messages
- `400` - invalid input params (friend)
- `405` - invalid method (not GET)

#### 5.1. GET Request

```bash
curl -i -X GET "http://localhost:9591/api/chan/size?friend=Bob"
```

#### 5.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Fri, 16 Jan 2026 18:56:08 GMT
Content-Length: 1

1
```

### 6. /api/chan/listen

- `200` - return message from broker
- `204` - no messages in consumer
- `400` - invalid input params (friend)
- `405` - invalid method (not GET)
- `406` - limit of registered subscribers
- `500` - invalid type of message container

#### 6.1. GET Request

```bash
curl -i -X GET "http://localhost:9591/api/chan/listen?friend=Bob"
```

#### 6.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 16 Jan 2026 18:57:22 GMT
Content-Length: 28

{"d":"hello","t":1768589838}
```
