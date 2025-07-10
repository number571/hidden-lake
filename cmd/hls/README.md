# HLS

> Hidden Lake Services

<img src="images/hls_logo.png" alt="hls_logo.png"/>

`Hidden Lake Services` ...

## List of services

1. [HLS=filesharer](hls_filesharer) - ...
2. [HLS=messenger](hls_messenger) - ...
3. [HLS=pinger](hls_pinger) - ...
4. [HLS=remoter](hls_remoter) - ...

## Installation

It is necessary to replace `<service>` with an existing service in the list.

```bash
$ go install github.com/number571/hidden-lake/cmd/hls/hls_<service>@latest
```

## How it works

...

<p align="center"><img src="images/hls_arch.png" alt="hls_arch.png"/></p>
<p align="center">Figure 1. Architecture of HLS</p>

## Supported platforms

- Windows (x86_64, arm64)
- Linux (x86_64, arm64)
- MacOS (x86_64, arm64)

## Build and run

Default build and run

```bash 
$ go run ./cmd/hls/hls_<service>

> [INFO] 2024/12/29 03:40:06 HLS=<service> is started
> ...
```

Creates [`./hls_<service>.yml`](./hls_messenger/hls_messenger.yml) file (as example `messenger`).

## Running options

```bash
$ hls_<service> --path /root
# path = path to config, database, key files
```
