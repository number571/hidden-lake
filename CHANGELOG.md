# CHANGELOG

<!-- ... -->

## v1.8.6~

*??? ??, ????*

### CHANGES

- `build/networks.yml`: rename settings, default_settings -> default, default_network
- `build/networks.yml`: delete queue_period, fetch_timeout
- `build/settings.yml`: replace many non global values to pkgs, services
- `cmd/hls,cmd/hla_tcp,cmd/hla_http`: message_size_bytes now optionally param (used default value)
- `modes`: rename dir configs -> modes
- `cmd/hlr`: add check password != ""

### BUG FIXES

- `cmd/hla/hla_http`: add check network mask
- `cmd/hlc`: adapters now not inserted in the config if HLC run with network key

<!-- ... -->

## v1.8.5

*May 27, 2025*

### IMPROVEMENTS

- `cmd/hla`: add hla_http service

### CHANGES

- `configs`: add classic, adaptive, multiplicative modes configs
- `build`: add overwrite default networks.yml, settings.yml => hl_networks.yml, hl_settings.yml

### BUG FIXES

- `cmd/hlm`: fix duplicate messages in WS receive

<!-- ... -->

## v1.8.4

*March 03, 2025*

### CHANGES

- `cmd/hls`: move threads param from args -> hls.yml config (pow_parallel, default=0->1)
- `build/settings.yml`: move consumers_cap -> hls.yml config (qbp_consumers, default=0->1)

<!-- ... -->

## v1.8.3

*January 11, 2025*

### IMPROVEMENTS

- `cmd/hlp`: new service - pinger

### CHANGES

- `internal/service/pkg/config`: rename GetLimitMessageSizeBytes -> GetPayloadSizeBytes
- `internal/utils/logger/anon`: addr, conn -> variable params
- `internal/applications/*`: rename config params incoming,interface -> external,internal
- `internal/applications/filesharer`: external param is now can be omitempty
- `internal/applications/messenger`: deleted ping function
- `internal/utils/pprof`: deleted

### BUG FIXES

- `cmd/hla/hla_tcp`: delete duplicate host add

<!-- ... -->

## v1.8.2

*December 21, 2024*

### CHANGES

- `cmd/*`: change args format input

<!-- ... -->

## v1.8.1

*December 15, 2024*

### CHANGES

- `pkg/*`: add custom errors
- `go.mod`: update go-peer version: v1.7.8 -> v1.7.9

### BUG FIXES

- `cmd/hls`: log broadcast WARN if len(connections) = 0

<!-- ... -->

## v1.8.0

*December 15, 2024*

### IMPROVEMENTS

- `cmd/hla`: rethink adapter's concept
- `cmd/hla/hla_tcp`: create

### CHANGES

- `cmd/*`: change log.Fatal -> panic (args validate)
- `cmd/*`: add *.yml default configs
- `cmd/hlt`: delete default config connection 127.0.0.1:9571
- `cmd/hls,cmd/hlc`: delete default config service hidden-lake-filesharer
- `internal/webui`: settings insert scheme://host and port -> insert scheme://host:port
- `cmd/hlt,cmd/hll,cmd/hle,cmd/hla/common`: deleted
- `go.mod`: update go-peer version: v1.7.6 -> v1.7.8

### BUG FIXES

- `cmd/hlm, cmd/hlf`: rename dir _daemon -> daemon

<!-- ... -->

## v1.7.7

*November 30, 2024*

### IMPROVEMENTS

- `cmd/hlm`: add support URL links

### CHANGES

- `cmd/*`: change 'GetXxxMS() uint64' methods to 'GetXxx() time.Duration' 
- `go.mod` [!]: update go-peer version: v1.7.3 -> v1.7.6
- `cmd/hlm,cmd/hlf`: move webui static, template paths to internal/webui
- `pkg/request,pkg/response`: update interfaces: IRequestBuilder, IResponseBuilder
- `hidden-lake`: move GVersion, GSettings, GNetworks from root dir -> to build/ dir
- `build`: default work_size_bits=22 -> work_size_bits=0
- `internal/utils/help`: create Println
- `internal/utils/name`: create IServiceName

### BUG FIXES

- `cmd/hle, cmd/hlt, cmd/hll`: fix serviceName in handlers
- `cmd/hlm`: fix bug downloadBase64File: filename can contains last char \

<!-- ... -->

## v1.7.6

*November 13, 2024*

### CHANGES

- `pkg`: add pkg/network
- `pkg`: move internal/service/pkg/request|response -> pkg/request|response
- `pkg/internal/utils/flag`: add key aliases
- `cmd/hls,cmd/hle`: rename parallel -> threads
- `cmd/*`: add 'help' arg
- `cmd/*/Dockerfile`: change SERVICE_PATH: "/mounted" -> "."

<!-- ... -->

## v1.7.5

*November 10, 2024*

### CHANGES

- `hidden-lake`: add build/settings.yml
- `hidden-lake`: move networks.yml -> build/networks.yml
- `networks.yml`: deleted j2BR39JfDf7Bajx3 network

<!-- ... -->

## v1.7.4

*November 05, 2024*

### IMPROVEMENTS

- `*`: test coverage > 80%

### CHANGES

- `cmd/hla/chatingar`: deleted
- `cmd/hls,cmd/hlc`: add sh daemon/checklast

### BUG FIXES

- `cmd/*`: move internal/config -> pkg/app/config

<!-- ... -->

## v1.7.3

*October 31, 2024*

### CHANGES

- `cmd/hlm,cmd/hlf`: add HLR,HLC,HLA links to HL services in /about page
- `cmd/hlm,cmd/hlf`: add target="_blank" to external links
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
