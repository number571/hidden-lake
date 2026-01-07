# HLS=remoter

> Hidden Lake Service (Remoter)

<img src="images/hls_remoter_logo.png" alt="hls_remoter_logo.png"/>

The `Hidden Lake Service (Remoter)` this is a service that provides the ability to make remote calls on the anonymous network core (HLK) with theoretically provable anonymity.

> [!CAUTION]
> This application can be extremely dangerous. Use HLS=remoter with caution.

> More information about HLS=remoter in the [habr.com/ru/articles/830130](https://habr.com/ru/articles/830130/ "Habr HLS=remoter")

## Installation

```bash
$ go install github.com/number571/hidden-lake/cmd/hls/hls-remoter@latest
```

## How it works

Most of the code is a call to API functions from the HLK kernel. Thanks to this approach, implicit authorization of users is formed from the state of the anonymizing service.

The server providing the remote access service is waiting for a request in the form of a command. The command does not depend on the operating system and therefore should have a small additional syntax separating the launch of the main command and its arguments.

As an example, to create a file with the contents of "hello, world!" and then reading from the same file, you will need to run the following command:

```bash
bash[@s]-c[@s]echo 'hello, world' > file.txt && cat file.txt
```

The `[@s]` label means that the arguments are separated for the main command.

## Supported platforms

- Windows (x86_64, arm64)
- Linux (x86_64, arm64)
- MacOS (x86_64, arm64)

## Build and run

Default build and run

```bash 
$ go run ./cmd/hls/hls-remoter

> [INFO] 2023/06/03 15:30:31 HLS=remoter is running...
> ...
```

Open port `9532` (HTTP, incoming).
Creates [`./hls-remoter.yml`](./hls-remoter.yml) file.

## Running options

```bash
$ hls-remoter --path /root
# path = path to config
```

## Example

The example will involve two nodes `recv_hlc, send_hls` and three repeaters `middle_hla_tcp_1, middle_hla_tcp_2, middle_hla_tcp_3`. The three remaining nodes are used only for the successful connection of the two main nodes. In other words, `HLA=tcp` nodes are traffic relay nodes.

Build and run nodes
```bash
$ cd examples/remoter/routing
$ make
```

Than run command
```bash
$ cd examples/remoter
$ make request # go run ./_request/main.go
```

Got response
```json
{"code":200,"head":{"Content-Type":"application/octet-stream"},"body":"aGVsbG8sIHdvcmxkCg=="}
```
