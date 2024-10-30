# CHANGELOG

<!-- ... -->

## v1.7.3~

*??? ??, ????*

### CHANGES

- `cmd/hlm,cmd/hlf`: add HLR,HLC,HLA links to HL services in /about page
- `cmd/hlf`: file hashing: sha256 -> sha384
- `cmd/hlm,cmd/hlf`: hash(pubkey.bytes()) -> hash(pubkey.string())

### BUG FIXES

- `cmd/hlm,cmd/hlf`: fix links to HL services in /about page
- `cmd/hlm`: fix emoji replacer

<!-- ... -->

## v1.7.2

*October 28, 2024*

### IMPROVEMENTS

- `cmd/hls,cmd/hle,cmd/hlt,cmd/hll`: add 'network' run option

### CHANGES

- `cmd/hls`: delete yaml host field from services in hls.yml
- `cmd/*`: delete default args from InitApp functions
- `cmd/hla/common`: simplified the code

### BUG FIXES

- `cmd/hle`: delete print log in decrypt failed block

<!-- ... -->

## v1.7.1

*October 24, 2024*

### CHANGES

- `cmd/hls`: rename /api/network/pubkey -> /api/service/pubkey
- `cmd/hle`: update API encrypt/decrypt messages
- `cmd/hls,cmd/hlt`: delete rand_ prefix parameters (message_size_bytes, queue_period_ms)
- `go.mod`: update go-peer version: 1.7.0 -> 1.7.2

<!-- ... -->

## v1.7.0

*October 20, 2024*

### IMPROVEMENTS

- `cmd/hlm`: add ping message
- `cmd`: RSA-4096 -> ML-KEM-768, ML-DSA-65

<!-- ... -->

## v1.6.21

*October 13, 2024*

### INIT
