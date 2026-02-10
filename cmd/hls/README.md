# HLS

> Hidden Lake Services

<img src="images/hls_logo.png" alt="hls_logo.png"/>

`Hidden Lake Services` are applied applications that perform specific tasks.

## List of services

1. [HLS=messenger](hls-messenger) - send and recv text messages
2. [HLS=filesharer](hls-filesharer) - view storage and download files 
3. [HLS=pinger](hls-pinger) - ping the node to check online status

## Installation

It is necessary to replace `<application>` with an existing service in the list.

```bash
$ go install github.com/number571/hidden-lake/cmd/hls/hls-<application>@latest
```

## How it works

HLS are classic http servers. The `Hlk-Sender-Name` header received from HLK is always present in requests. The service may return the `Hlk-Response-Mode` header, indicating that the response must be sent (on/off). Services also have a limit on the size of messages sent = `payload_size_bytes` (HLK).

## Supported platforms

- Windows (x86_64, arm64)
- Linux (x86_64, arm64)
- MacOS (x86_64, arm64)

## Build and run

Default build and run

```bash 
$ go run ./cmd/hls/hls-<application>

> [INFO] 2024/12/29 03:40:06 HLS=<application> is started
> ...
```

Creates [`./hls-<application>.yml`](./hls-messenger/hls-messenger.yml) file (as example `messenger`). Also can create another files / directories (as example `hls-filesharer.stg` in `filesharer`)

## Running options

```bash
$ hls-<application> --path /root
# path    = path to config, database, key files
```
