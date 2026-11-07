# HLS=notifier

> Hidden Lake Service (Notifier)

<img src="images/hls_notifier_logo.png" alt="hls_notifier_logo.png"/>

The `Hidden Lake Service (Notifier)` is a notifier based on the core of an anonymous network with theoretically provable anonymity of HLK. A feature of this notifier is the provision of anonymity of the fact of transactions (sending, receiving) bytes messages.

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
Creates [`./hls-notifier.yml`](./hls-notifier.yml) file.

## Running options

```bash
$ hls-notifier --path /root
# path = path to config file
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
// RECOMMENDATION: empty body is equivalent to a PING request
[uint32]{"method":"POST","host":"hls-notifier","path":"/notify"}[body]
[body]=<bytes>

// header response size = (4 + 12) = 16B
// body response size can be = (144 - 16) = 128B
[uint32]{"code":200}[body]
```

## HLS API

```
1. GET      /api/index          | params = [] 
                                |> description = get name of service
2. GET/POST /api/chan/notify    | params = ["friend":string]
                                |> description = get limit notify size / send notify to chan
3. GET      /api/chan/listen    | params = ["sid":string]
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
curl -i -X POST "http://localhost:9591/api/chan/notify?friend=Bob" --data 'hello, world!'
```

#### 2.2. POST Response

```
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Fri, 16 Jan 2026 18:55:12 GMT
Content-Length: 2

ok
```

### 3. /api/chan/listen

- `200` - return message from broker
- `204` - no messages in consumer
- `400` - invalid input params (friend)
- `405` - invalid method (not GET)
- `406` - limit of registered subscribers
- `500` - invalid type of message container

#### 3.1. GET Request

```bash
curl -i -X GET "http://localhost:9591/api/chan/listen?sid=some_application_id"
```

#### 3.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Fri, 16 Jan 2026 18:57:22 GMT
Content-Length: 13

hello, world!
```
