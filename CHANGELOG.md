# CHANGELOG

<!-- ... -->

## v1.10.2~

*??? ??, ????*

### IMPROVEMENTS

- `cmd/hls/hls-filesharer`: add API function DelRemoteFile

### CHANGES

- `cmd/hls/hls-remoter`: deleted
- `cmd/hls/hls-filesharer`: change API function GetRemoteFile
- `cmd/hls/hls-filesharer`: The file names in private are now different when downloading a file with the same name from the personal and sharing directories
- `pkg/api/*/client`: change GetIndex functions
- `cmd/hlc`: add hls-filesharer to init config

<!-- ... -->

## v1.10.1

*February 05, 2026*

### IMPROVEMENTS

- `cmd/hls/hls-filesharer`: add check hash/size of file on load chunk

### CHANGES

- `cmd/hls/*`: move requests to HLK from HLS -> pkg/api/services/*/request

### BUG FIXES

- `cmd/hls/hls-filesharer`: fix check error on loadChunk (status code)

<!-- ... -->

## v1.10.0

*February 02, 2026*

### IMPROVEMENTS

- `*/pkg/app/init_app.go`: create path for configs/databases if not exists
- `cmd/hls/*|cmd/hlk`: add CLI clients
- `cmd/hls/*`: add internal API, update pkg/client, update examples
- `cmd/hls/hls-filesharer`: add info API
- `*/pkg/client|*/pkg/config`: move from internal -> pkg/api
- `cmd/hls/hls-filesharer`: switch global storage to (private, sharing) with personal

### CHANGES

- `go.mod`: update go-peer
- `*/daemon/install_*.sh`: change path /root -> (/usr/local/bin, ~/.config/hidden-lake)
- `cmd/hls/hls-filesharer|cmd/hls/hls-messenger`: delete webui -> change to cli
- `cmd/hls/hls-remoter`: change separator [@remoter-separator] -> [@s]
- `cmd/hlk`: change API of the /api/network/request path
- `pkg`: change path /request, /response, /handler, /adapters -> /network/request, /network/response, /network/handler, /network/adapters
- `build`: move pkg/build -> build/environment
- `*`: rename structs S/name/Error -> SError
- `cmd/hlk`: add sort list of friends (get function)
- `cmd/hls-messenger`: change websocket -> longpoll method on listen messages from chat

### BUG FIXES

- `cmd/hla/hla-tcp`: fix duplicated log of start app
- `pkg/adapters/http/client`: add missed GetSettings
- `cmd/hls/hls-filesharer`: delete header CHeaderResponseMode from calculate size of response message

<!-- ... -->

## v1.9.1

*October 05, 2025*

### CHANGES

- `cmd/hlk`: change API '/api/kernel/pubkey' -> '/api/profile/pubkey'
- `cmd/hlk`: change headers CHeaderSenderName, CHeaderResponseMode
- `cmd/hlc,cmd/hlk`: change applications/services names

### BUG FIXES

- `internal/adapters/http`: fix get_connections scheme
- `cmd/hls/hls-filesharer`: fix stream resume download

<!-- ... -->

## v1.9.0

*August 17, 2025*

### CHANGES

- `*_test.go`: replace {t.Error(); return} -> {t.Fatal()}
- `cmd/hla/hla-http`: delete write_timeout_ms
- `cmd`: renames hls -> hlk; hlm, hlf, hlr, hlp -> hls-messenger, hls-filesharer, hls-remoter, hls-pinger
- `cmd`: rename hla_tcp -> hla-tcp, hla_http -> hla-http, hls_messenger -> hls-messenger, ...

### BUG FIXES

- `cmd/*`: fix print information about version, help
- `cmd/hla/hla-http`: delete request_timeout_ms

<!-- ... -->

## v1.8.6

*June 13, 2025*

### IMPROVEMENTS

- `test/prod`: add test prod env

### CHANGES

- `build/networks.yml`: rename settings, default_settings -> default, default_network
- `build/networks.yml`: delete queue_period, fetch_timeout
- `build/settings.yml`: replace many non global values to pkgs, services, add http timeouts
- `cmd/hls,cmd/hla-tcp,cmd/hla-http`: message_size_bytes now optionally param (used default value)
- `modes`: rename dir configs -> modes
- `cmd/hls/hls-remoter`: add check password != ""

### BUG FIXES

- `cmd/hla/hla-http`: add check network mask
- `cmd/hlc`: adapters now not inserted in the config if HLC run with network key

<!-- ... -->

## v1.8.5

*May 27, 2025*

### IMPROVEMENTS

- `cmd/hla`: add hla-http service

### CHANGES

- `configs`: add classic, adaptive, multiplicative modes configs
- `build`: add overwrite default networks.yml, settings.yml => hl-networks.yml, hl-settings.yml

### BUG FIXES

- `cmd/hls/hls-messenger`: fix duplicate messages in WS receive

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

- `cmd/hls/hls-pinger`: new service - pinger

### CHANGES

- `internal/service/pkg/config`: rename GetLimitMessageSizeBytes -> GetPayloadSizeBytes
- `internal/utils/logger/anon`: addr, conn -> variable params
- `internal/applications/*`: rename config params incoming,interface -> external,internal
- `internal/applications/filesharer`: external param is now can be omitempty
- `internal/applications/messenger`: deleted ping function
- `internal/utils/pprof`: deleted

### BUG FIXES

- `cmd/hla/hla-tcp`: delete duplicate host add

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
- `cmd/hla/hla-tcp`: create

### CHANGES

- `cmd/*`: change log.Fatal -> panic (args validate)
- `cmd/*`: add *.yml default configs
- `cmd/hlt`: delete default config connection 127.0.0.1:9571
- `cmd/hls,cmd/hlc`: delete default config service hls-filesharer
- `internal/webui`: settings insert scheme://host and port -> insert scheme://host:port
- `cmd/hlt,cmd/hll,cmd/hle,cmd/hla/common`: deleted
- `go.mod`: update go-peer version: v1.7.6 -> v1.7.8

### BUG FIXES

- `cmd/hls/hls-messenger, cmd/hls/hls-filesharer`: rename dir _daemon -> daemon

<!-- ... -->

## v1.7.7

*November 30, 2024*

### IMPROVEMENTS

- `cmd/hls/hls-messenger`: add support URL links

### CHANGES

- `cmd/*`: change 'GetXxxMS() uint64' methods to 'GetXxx() time.Duration' 
- `go.mod` [!]: update go-peer version: v1.7.3 -> v1.7.6
- `cmd/hls/hls-messenger,cmd/hls/hls-filesharer`: move webui static, template paths to internal/webui
- `pkg/request,pkg/response`: update interfaces: IRequestBuilder, IResponseBuilder
- `hidden-lake`: move GVersion, GSettings, GNetworks from root dir -> to build/ dir
- `build`: default work_size_bits=22 -> work_size_bits=0
- `internal/utils/help`: create Println
- `internal/utils/name`: create IServiceName

### BUG FIXES

- `cmd/hle, cmd/hlt, cmd/hll`: fix serviceName in handlers
- `cmd/hls/hls-messenger`: fix bug downloadBase64File: filename can contains last char \

<!-- ... -->

## v1.7.6

*November 13, 2024*

### CHANGES

- `pkg`: add pkg/network
- `pkg`: move internal/service/pkg/request|response -> pkg/request|response
- `pkg/pkg/utils/flag`: add key aliases
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

- `cmd/hls/hls-messenger,cmd/hls/hls-filesharer`: add HLS=remoter,HLC,HLA links to HL services in /about page
- `cmd/hls/hls-messenger,cmd/hls/hls-filesharer`: add target="_blank" to external links
- `cmd/hls/hls-filesharer`: file hashing: sha256 -> sha384
- `cmd/hls/hls-messenger,cmd/hls/hls-filesharer`: hash(pubkey.bytes()) -> hash(pubkey.string())

### BUG FIXES

- `cmd/hls/hls-messenger,cmd/hls/hls-filesharer`: fix links to HL services in /about page
- `cmd/hls/hls-messenger`: fix emoji replacer

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

- `cmd/hls`: rename /api/network/pubkey -> /api/profile/pubkey
- `cmd/hle`: update API encrypt/decrypt messages
- `cmd/hls,cmd/hlt`: delete rand_ prefix parameters (message_size_bytes, queue_period_ms)
- `go.mod`: update go-peer version: 1.7.0 -> 1.7.2

<!-- ... -->

## v1.7.0

*October 20, 2024*

### IMPROVEMENTS

- `cmd/hls/hls-messenger`: add ping message
- `cmd`: RSA-4096 -> ML-KEM-768, ML-DSA-65

<!-- ... -->

## v1.6.21

*October 13, 2024*

### INIT
