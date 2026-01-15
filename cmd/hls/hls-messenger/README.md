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

Than run command
```bash
$ cd examples/messenger
$ ./_request/raw/request.sh
```

Got response
```
success: send
```
