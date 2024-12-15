# HLC

> Hidden Lake Composite

<img src="images/hlc_logo.png" alt="hlc_logo.png"/>

The `Hidden Lake Composite` combines several HL type's services into one application using startup config.

## Installation

```bash
$ go install github.com/number571/hidden-lake/cmd/hlc@latest
```

## How it works

The application HLC includes the download of all Hidden Lake services, and runs only the configurations selected by names in the file. The exact names of the services can be found in their `pkg/settings/settings.go` configuration files.

## Supported platforms

- Windows (x86_64, arm64)
- Linux (x86_64, arm64)
- MacOS (x86_64, arm64)

## Build and run

Default build and run

```bash 
$ go run ./cmd/hlc

> [INFO] 2023/12/03 02:12:51 HLC is running...
> ...
```

Creates [`./hlc.yml`](./hlc.yml) file.

## Running options

```bash
$ hlc -path=/root -network=xxx -threads=1
# path    = path to config, database, key files
# network = use network configuration from networks.yml
# threads = num of parallel functions for PoW algorithm
```

## Config structure

```
"logging"  Enable loggins in/out actions in the network
"services" Names of Hidden Lake services 
```
