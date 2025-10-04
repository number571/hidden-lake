# HLA=http

> Hidden Lake Adapter (HTTP)

<img src="images/hla_http_logo.png" alt="hla_http_logo.png"/>

The `Hidden Lake Adapter (HTTP)` allows adapt HL traffic based on the HTT[] protocol.

## Installation

```bash
$ go install github.com/number571/hidden-lake/cmd/hla/hla-http@latest
```

## How it works

HLA=http uses `/api/network/adapter` handle function to Produce/Consume HL messages.

## Supported platforms

- Windows (x86_64, arm64)
- Linux (x86_64, arm64)
- MacOS (x86_64, arm64)

## Build and run

Default build and run

```bash 
$ go run ./cmd/hla/hla-http

> [INFO] 2023/06/03 15:30:31 HLA=http is running...
> ...
```

Open port `9511` (HTTP, external), `9512` (HTTP, internal).
Creates [`./hla-http.yml`](./hla-http.yml) file.

## Running options

```bash
$ hla-http --path /root
# path = path to config
```
