# HLS=messenger-CLI

## Installation

```bash
$ go install github.com/number571/hidden-lake/cmd/hls/hls-messenger/client/hls-messenger-cli@latest
```

## Running options

```bash
$ hls-messenger-cli --service {{HLS-address}} --friend {{friend-name}}
# service = address of the HLS internal address (default localhost:9591)
# friend  = alias name of the friend
```

> You can get more detailed information using the `--help` option.

### Examples

```bash
$ hls-messenger-cli --friend {{friend-name}}
{
        "friend_name": "{{friend-name}}",
        "payload_limit": 3552
}

hello, world!
```
