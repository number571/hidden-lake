# HLS=remoter-CLI

## Installation

```bash
$ go install github.com/number571/hidden-lake/cmd/hls/hls-remoter/client/hls-remoter-cli@latest
```

## Running options

```bash
$ hls-remoter-cli --service {{HLS-address}} --friend {{friend-name}}
# service = address of the HLS internal address (default localhost:9531)
# friend  = alias name of the friend
```

> You can get more detailed information using the `--help` option.

### Examples

```bash
$ hls-remoter-cli --friend {{friend-name}}
password: {{remote-password}}
> ls
```
