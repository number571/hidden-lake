# HLH=keyops

> Hidden Lake Helper (Keyops)

<img src="images/hlh_keyops_logo.png" alt="hlh_keyops_logo.png"/>

The `Hidden Lake Helper (Keyops)` application provides generate / encapsulate / decapsulate key operations.

## Installation

```bash
$ go install github.com/number571/hidden-lake/cmd/hlh/hlh-keyops@latest
```

## How it works

The application HLH=keyops uses `go-peer/pkg/crypto/asymmetric` logic.

## Supported platforms

- Windows (x86_64, arm64)
- Linux (x86_64, arm64)
- MacOS (x86_64, arm64)

## Build and run

Default build and run

```bash 
$ go run ./cmd/hlh/hlh-keyops

<Hidden Lake Helper = Keyops (HLH=keyops)>
...
```

## Examples

### 1. Generate key pair

```bash
$ hlh-keyops -p "/path/to/my_keypair" -g
```

### 2. Encapsulate key

```bash
$ hlh-keyops -e "/path/to/public_key" -o "/path/to/encapsulated_key"
> <STDOUT: generated secret key in the hex encode>
```

### 3. Decapsulate key

```bash
$ hlh-keyops -p "/path/to/my_keypair" -d "/path/to/encapsulated_key"
> <STDOUT: decapsulated secret key in the hex encode>
```
