# HLS

> Hidden Lake Services

<img src="images/hls_logo.png" alt="hls_logo.png"/>

`Hidden Lake Services` are applied applications that perform specific tasks.

## List of services

1. [HLS=filesharer](hls-filesharer) - file sharing with a web interface
2. [HLS=messenger](hls-messenger) - chat with a web interface
3. [HLS=pinger](hls-pinger) - ping the node to check the online status
4. [HLS=remoter](hls-remoter) - executes remote access commands

## Installation

It is necessary to replace `<application>` with an existing service in the list.

```bash
$ go install github.com/number571/hidden-lake/cmd/hls/hls_<application>@latest
```

## How it works

HLS are classic http servers. The `Hlk-Sender-Pubkey` header received from HLK is always present in requests. The service may return the `Hlk-Response-Mode` header, indicating that the response must be sent (on/off). Services also have a limit on the size of messages sent = `payload_size_bytes` (HLK).

## Supported platforms

- Windows (x86_64, arm64)
- Linux (x86_64, arm64)
- MacOS (x86_64, arm64)

## Build and run

Default build and run

```bash 
$ go run ./cmd/hls/hls_<application>

> [INFO] 2024/12/29 03:40:06 HLS=<application> is started
> ...
```

Creates [`./hls_<application>.yml`](./hls-messenger/hls-messenger.yml) file (as example `messenger`).

## Running options

```bash
$ hls_<application> --path /root
# path    = path to config, database, key files
```
