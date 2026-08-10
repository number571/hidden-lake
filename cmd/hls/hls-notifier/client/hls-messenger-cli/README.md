# HLS=notifier-CLI

## Installation

```bash
$ go install github.com/number571/hidden-lake/cmd/hls/hls-notifier/client/hls-notifier-cli@latest
```

## Running options

```bash
$ hls-notifier-cli --service {{HLS-address}} --friend {{friend-name}}
# service = address of the HLS internal address (default localhost:9591)
# friend  = alias name of the friend
```

> You can get more detailed information using the `--help` option.

### Examples

```bash
$ hls-notifier-cli --friend {{friend-name}}
{
        "friend_name": "{{friend-name}}",
        "payload_limit": 3552
}

hello, world!
```
