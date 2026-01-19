# HLS=pinger-CLI

## Installation

```bash
$ go install github.com/number571/hidden-lake/cmd/hls/hls-pinger/client/hls-pinger-cli@latest
```

## Running options

```bash
$ hls-pinger-cli --service {{HLS-address}} --friend {{friend-name}}
# service = address of the HLS internal address (default localhost:9551)
# friend  = alias name of the friend
```

> You can get more detailed information using the `--help` option.

### Examples

```bash
$ hls-pinger-cli --friend {{friend-name}}
ok
```
