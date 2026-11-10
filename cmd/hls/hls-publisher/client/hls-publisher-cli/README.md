# HLS=publisher-CLI

## Installation

```bash
$ go install github.com/number571/hidden-lake/cmd/hls/hls-publisher/client/hls-publisher-cli@latest
```

## Running options

```bash
$ hls-publisher-cli --service {{HLS-address}} --friend {{friend-name}} --do {{function-name}} --arg {{arguments}} --remote --public
# service = address of the HLS internal address (default localhost:9551)
# friend  = alias name of the friend
# do      = function to run [size|push|load-by-uid|load-by-nid|listen]
# arg     = arguments of the function
# remote  = send/recv messages to/from node storage (default false)
# public  = use public method of send/recv messages (default false)
```

> You can get more detailed information using the `--help` option.

### Examples

1. Get count of messages

```bash
### Get count of messages from local box (personal) ###
$ hls-publisher-cli -f {{friend-name}} -d size
```

```bash
### Get count of messages from local box (public) ###
$ hls-publisher-cli -f {{friend-name}} -d size --public

```bash
### Get count of messages from friend's box (personal from node) ###
$ hls-publisher-cli -f {{friend-name}} -d size --remote
```

```bash
### Get count of messages from friend's box (public from node) ###
$ hls-publisher-cli -f {{friend-name}} -d size --remote --public
```

2. Push message

```bash
### Push message to local friend's box (personal) ###
$ hls-publisher-cli -f {{friend-name}} -d push -a "{{message}}"
```

```bash
### Push message to remote friend's box (personal to node) ###
$ hls-publisher-cli -f {{friend-name}} -d push -a "{{message}}" --remote
```

```bash
### Push message to local public box (public) ###
$ hls-publisher-cli -f {{friend-name}} -d push -a "{{message}}" --public
```

```bash
# (public to node not allowed!)
```

3. Listen messages

```bash
### Subcribe on listen new (personal) messages from nodes ###
$ hls-publisher-cli -f {{friend-name}} -d listen
```

4. Load message by UID (uniq ID)

```bash
### Get message from local box (personal) ###
$ hls-publisher-cli -f {{friend-name}} -d load-by-uid -a "{{UID}}"
```

```bash
### Get message from remote box (personal from node) ###
$ hls-publisher-cli -f {{friend-name}} -d load-by-uid -a "{{UID}}" --remote
```

```bash
### Get message from local box (public) ###
$ hls-publisher-cli -f {{friend-name}} -d load-by-uid -a "{{UID}}" --public
```

```bash
### Get message from remote box (public from node) ###
$ hls-publisher-cli -f {{friend-name}} -d load-by-uid -a "{{UID}}" --remote --public
```

5. Load message by NID (number ID)

```bash
### Get message from local box (personal) ###
$ hls-publisher-cli -f {{friend-name}} -d load-by-nid -a "{{NID}}"
```

```bash
### Get message from remote box (personal from node) ###
$ hls-publisher-cli -f {{friend-name}} -d load-by-nid -a "{{NID}}" --remote
```

```bash
### Get message from local box (public) ###
$ hls-publisher-cli -f {{friend-name}} -d load-by-nid -a "{{NID}}" --public
```

```bash
### Get message from remote box (public from node) ###
$ hls-publisher-cli -f {{friend-name}} -d load-by-nid -a "{{NID}}" --remote --public
```
