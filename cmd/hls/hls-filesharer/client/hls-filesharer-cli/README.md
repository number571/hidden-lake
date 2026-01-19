# HLS=filesharer-CLI

## Installation

```bash
$ go install github.com/number571/hidden-lake/cmd/hls/hls-filesharer/client/hls-filesharer-cli@latest
```

## Running options

```bash
$ hls-filesharer-cli --service {{HLS-address}} --path {{path-to-files}} --friend {{friend-name}} --do {{function-name}} --arg {{arguments}}
# service = address of the HLS internal address (default localhost:9541)
# path    = path to store download files  (default .)
# friend  = alias name of the friend
# do      = function to run
# arg     = arguments of the function
```

> You can get more detailed information using the `--help` option.

### Examples

1. List

```bash
### Get list of files from friend ###
$ hls-filesharer-cli -f {{friend-name}} -d list -a {{page}}
```

2. Info

```bash
### Get info of the file from friend ###
$ hls-filesharer-cli -f {{friend-name}} -d info -a {{filename}}
```

3. Load

```bash
### Download file from friend ###
$ hls-filesharer-cli -f {{friend-name}} -d load -a {{filename}}
```
