# HLA=meshtastic

> Hidden Lake Adapter (Meshtastic)

<img src="images/hla_meshtastic_logo.png" alt="hla_meshtastic_logo.png"/>

The `Hidden Lake Adapter (Meshtastic)` allows adapt HL traffic based on the Meshtastic/LoRa protocol.

## Installation

```bash
$ go install github.com/number571/hidden-lake/cmd/hla/hla-meshtastic@latest
```

## Requirements

1. Python version `== 3.14.5`
2. Connected Meshtastic/LoRa device via USB

## How it works

HLA=meshtastic uses Python [script.py](../../../pkg/network/adapters/meshtastic/service/script.py) for send/recv messages. Before launching, dependencies for Python are automatically loaded from file [requirements.txt](../../../pkg/network/adapters/meshtastic/service/requirements.txt).

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
$ hla-meshtastic --path /root --network xxx
# path    = path to config
# network = use network configuration from networks.yml
```
