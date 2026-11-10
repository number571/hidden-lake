# HLS=publisher

> Hidden Lake Service (Publisher)

<img src="images/hls_publisher_logo.png" alt="hls_publisher_logo.png"/>

The `Hidden Lake Service (Publisher)` is a publisher based on the core of an anonymous network with theoretically provable anonymity of HLK. A feature of this publisher is the provision of anonymity of the fact of transactions (pushing, loading) text messages.

> More information about HLS=publisher in the [habr.com/ru/post/701488](https://habr.com/ru/post/701488/ "Habr HLS=publisher")

## Installation

```bash
$ go install github.com/number571/hidden-lake/cmd/hls/hls-publisher@latest
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
$ go run ./cmd/hls/hls-publisher

> [INFO] 2023/06/03 15:30:31 HLS=publisher is running...
> ...
```

Open ports `9551` (HTTP, internal) and `9552` (HTTP, incoming).
Creates [`./hls-publisher.yml`](./hls-publisher.yml) and `./hls-publisher.db` files.

## Running options

```bash
$ hls-publisher --path /root
# path = path to config and database files
```

## Example

The example will involve (as well as in HLK) five nodes `node1_hlm, node2_hlm` and `middle_hla_tcp_1, middle_hla_tcp_2, middle_hla_tcp_3`. The three `HLA=tcp` nodes are only needed for communication between `node1_hlm` and `node2_hlm` nodes. Each of the remaining ones is a combination of HLK and HLS=publisher, where HLS=publisher plays the role of an application and services (as it was depicted in `Figure 3` HLK readme).

Build and run nodes
```bash
$ cd examples/publisher/routing
$ make
```

Than run commands
```bash
### Terminal 1 ###
$ cd examples/publisher
$ make request-test-listen
### Terminal 2 ###
$ cd examples/publisher
$ ARG='"hello, world!"' make request-test-push-remote
```

## HLS/HLK API

### 1. /push

- `200` - return index
- `202` - message already exist in database
- `400` - invalid input params (uid, body)
- `404` - message by uid not found
- `405` - invalid method (not GET)
- `409` - failed read HTTP request body
- `500` - load/push from database
- `502` - get aliasName HLK

```json
// header request size = (4 + 82) = 86B
// body request size can be = (144 - 86) = 58B
// INFORMATION: uid = random chars from [go-peer/pkg/crypto/random/random.go]
[uint32]{"method":"POST","host":"hls-publisher","path":"/push?uid=0000000000000000000000"}[body]
[body]=<bytes>

// header response size = (4 + 12) = 16B
// body response size can be = (144 - 16) = 128B
[uint32]{"code":200}[body]
[body]=<uint64>
```

### 2. /load

```json
// header request size = (4 + 91) = 95B
// body request size can be = (144 - 95) = 49B
[uint32]{"method":"GET","host":"hls-publisher","path":"/load?uid=0000000000000000000000&pub=false"}[body]
// | OR
// header request size = (4 + 89) = 93B
// body request size can be = (144 - 93) = 51B
[uint32]{"method":"GET","host":"hls-publisher","path":"/load?nid=18446744073709551616&pub=false"}[body]

// header response size = (4 + 60) = 64B
// body response size can be = (144 - 64) = 80B
[uint32]{"code":200,"head":{"M-S":"1","M-I":"18446744073709551616"}}[body]
[body]=<bytes>
// | OR
// header response size = (4 + 62) = 66B
// body response size can be = (144 - 66) = 78B
[uint32]{"code":200,"head":{"M-S":"1","M-I":"0000000000000000000000"}}[body]
[body]=<bytes>
```

### 3. /size

```json
// header request size = (4 + 64) = 68B
// body request size can be = (144 - 68) = 76B
[uint32]{"method":"GET","host":"hls-publisher","path":"/size&pub=false"}[body]

// header response size = (4 + 12) = 16B
// body response size can be = (144 - 16) = 128B
[uint32]{"code":200}[body]
[body]=<uint64>
```

## HLS API

```
1. GET          /api/index          | params = [] 
                                    |> description = get name of service
3. GET/POST     /api/msgbox/push    | params = ["friend":string,"remote":?uint8,"public":?bool]
                                    |> description = get limit message size / send message to box
4. GET          /api/msgbox/load    | params = ["friend":string,"nid":?uint64,"uid":?string,"remote":?uint8,"public":?bool]
                                    |> description = get message from box by index / uid
5. GET          /api/msgbox/size    | params = ["friend":string,"remote":?uint8,"public":?bool]
                                    |> description = get count of messages in the box
6. GET          /api/msgbox/listen  | params = ["sid":string]
                                    |> description = try get message from friends with longpoll method 
```

### 1. /api/index

- `200` - return "hidden-lake-service=publisher"

#### 1.1. GET Request

```bash
curl -i -X GET http://localhost:9551/api/index
```

#### 1.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Fri, 16 Jan 2026 18:54:07 GMT
Content-Length: 29

hidden-lake-service=publisher
```

### 2. /api/msgbox/push

- `200` - return limit size (GET) or timestamp (POST)
- `400` - invalid input params (friend, body)
- `405` - invalid method (not GET, POST)
- `409` - failed read HTTP request body
- `502` - failed send request (HLK side)

#### 2.1. GET Request

```bash
curl -i -X GET "http://localhost:9551/api/chan/notify"
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
curl -i -X POST "http://localhost:9551/api/chan/notify?friend=Bob&uid=0000000000000000000000" --data 'hello, world!'
```

#### 2.2. POST Response

```
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Fri, 16 Jan 2026 18:55:12 GMT
Content-Length: 19

2026-01-16T18:55:12
```

### 3. /api/chan/msgbox/push

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
curl -i -X GET "http://localhost:9551/api/chan/msgbox/push"
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
curl -i -X POST "http://localhost:9551/api/chan/msgbox/push?friend=Bob&uid=0000000000000000000000" --data 'hello, world!'
```

#### 3.2. POST Response

```
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Fri, 16 Jan 2026 18:55:12 GMT
Content-Length: 19

2026-01-16T18:55:12
```

### 4. /api/chan/msgbox/load

- `200` - return message from database / node
- `400` - invalid input params (friend, index)
- `404` - message not found in database
- `405` - invalid method (not GET)
- `418` - got bad response from node
- `500` - load message from database
- `503` - request was sent but no response was received

#### 4.1. GET Request

```bash
curl -i -X GET "http://localhost:9551/api/chan/msgbox/load?friend=Bob&nid=0"
```

#### 4.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 16 Jan 2026 18:56:08 GMT
Content-Length: 89

{"head":{"s":2,"u":"0000000000000000000000"},"body":{"d":"hello, world!","t":1768589712}}
```

### 5. /api/chan/msgbox/size

- `200` - return count of messages
- `400` - invalid input params (friend)
- `405` - invalid method (not GET)

#### 5.1. GET Request

```bash
curl -i -X GET "http://localhost:9551/api/chan/msgbox/size?friend=Bob"
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
curl -i -X GET "http://localhost:9551/api/chan/listen?friend=Bob"
```

#### 6.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 16 Jan 2026 18:57:22 GMT
Content-Length: 28

{"d":"hello","t":1768589838}
```
