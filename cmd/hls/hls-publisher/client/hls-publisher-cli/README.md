# HLS=publisher-CLI

## Installation

```bash
$ go install github.com/number571/hidden-lake/cmd/hls/hls-publisher/client/hls-publisher-cli@latest
```

## Running options

```bash
$ hls-publisher-cli --service {{HLS-address}} --friend {{friend-name}}
# service = address of the HLS internal address (default localhost:9551)
# friend  = alias name of the friend
```

> You can get more detailed information using the `--help` option.

### Examples

```bash
$ hls-publisher-cli --friend {{friend-name}}
{
        "friend_name": "{{friend-name}}",
        "payload_limit": 3552
}

hello, world!
```
