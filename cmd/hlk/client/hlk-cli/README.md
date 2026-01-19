# HLK-CLI

## Installation

```bash
$ go install github.com/number571/hidden-lake/cmd/hlk/client/hlk-cli@latest
```

## Running options

```bash
$ hlk-cli --kernel {{HLK-address}} --do {{function-name}} --arg {{arguments}}
# kernel = address of the HLK internal address (default localhost:9572)
# do     = function to run
# arg    = arguments of the function
```

> You can get more detailed information using the `--help` option.

### Examples

1. Network

```bash
### Send message without response (example from hls-messenger) ###
$ echo '{"method":"POST","host":"hls-messenger","path":"/push","body":"aGVsbG8sIHdvcmxkIQ=="}' | hlk-cli -d send -a {{friend-name}}
```

```bash
### Send message and get response (example from hls-pinger) ###
$ echo '{"method":"GET","host":"hls-pinger","path":"/ping"}' | hlk-cli -d fetch -a {{friend-name}}
```

2. Profile

```bash
### Get own public key ###
$ hlk-cli -d pubkey
```

3. Onlines

```bash
### Get alive connections ###
$ hlk-cli -d get-onlines
```

```bash
### Delete alive connection ###
$ hlk-cli -d del-online -a tcp://localhost:9999
```

4. Connections

```bash
### Get all connections ###
$ hlk-cli -d get-connections
```

```bash
### Add connection ###
$ hlk-cli -d add-connection -a tcp://localhost:9999
```

```bash
### Delete connection ###
$ hlk-cli -d del-connection -a tcp://localhost:9999
```

5. Friends

```bash
### Get all friends ###
$ hlk-cli -d get-friends
```

```bash
### Add friend ###
$ echo "PubKey{...}" | hlk-cli -d add-friend -a {{friend-name}}
```

```bash
### Delete friend ###
$ hlk-cli -d del-friend -a {{friend-name}}
```
