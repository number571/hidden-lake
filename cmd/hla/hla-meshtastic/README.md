# HLA=meshtastic

> Hidden Lake Adapter (TCP)

<img src="images/hla_meshtastic_logo.png" alt="hla_meshtastic_logo.png"/>

The `Hidden Lake Adapter (Meshtastic)` allows adapt HL traffic based on the Meshtastic/LoRa protocol.

## Installation

```bash
$ go install github.com/number571/hidden-lake/cmd/hla/hla-meshtastic@latest
```

## How it works

HLA=meshtastic uses python [script](../../../pkg/network/adapters/meshtastic/service/script.py).

## Supported platforms

- Windows (x86_64, arm64)
- Linux (x86_64, arm64)
- MacOS (x86_64, arm64)

## Build and run

Default build and run

```bash 
$ go run ./cmd/hla/hla-meshtastic

> [INFO] 2023/06/03 15:30:31 HLA=meshtastic is running...
> ...
```

Open port `9501` (HTTP, external), `9502` (HTTP, internal).
Creates [`./hla-meshtastic.yml`](./hla-meshtastic.yml) file.

## Running options

```bash
$ hla-tmeshtasticcp --path /root --network xxx
# path    = path to config
# network = use network configuration from networks.yml
```
