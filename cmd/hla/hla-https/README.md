# HLA=https

> Hidden Lake Adapter (HTTPS)

<img src="images/hla_https_logo.png" alt="hla_https_logo.png"/>

The `Hidden Lake Adapter (HTTPS)` allows adapt HL traffic based on the HTTPS protocol.

## Installation

```bash
$ go install github.com/number571/hidden-lake/cmd/hla/hla-https@latest
```

## How it works

HLA=https uses `/adapter/produce`, `/adapter/consume` handle functions to Produce/Consume HL messages.

## Supported platforms

- Windows (x86_64, arm64)
- Linux (x86_64, arm64)
- MacOS (x86_64, arm64)

## Build and run

Default build and run

```bash 
$ go run ./cmd/hla/hla-https

> [INFO] 2023/06/03 15:30:31 HLA=https is running...
> ...
```

Open port `9531` (HTTPS, external), `9532` (HTTPS, internal).
Creates [`./hla-https.yml`](./hla-https.yml) file.

## Running options

```bash
$ hla-https --path /root --network xxx
# path    = path to config
# network = use network configuration from networks.yml
```
