# HLS=filesharer-CLI

## Installation

```bash
$ go install github.com/number571/hidden-lake/cmd/hls/hls-filesharer/client/hls-filesharer-cli@latest
```

## Running options

```bash
$ hls-filesharer-cli --service {{HLS-address}} --path {{path-to-files}} --type {{type-of-storage}} --friend {{friend-name}} --do {{function-name}} --arg {{arguments}}
# service = address of the HLS internal address (default localhost:9541)
# path    = path to store download files (default .)
# type    = type of storage [local|personal|public] (default public)
# friend  = alias name of the friend
# do      = function to run [list|info|download|upload|delete]
# arg     = arguments of the function
```

> You can get more detailed information using the `--help` option.

### Examples

1. Remote list

```bash
### Get list of files from friend (public) ###
$ hls-filesharer-cli -f {{friend-name}} -d list -a {{page}}
```

```bash
### Get list of files from friend (personal) ###
$ hls-filesharer-cli -f {{friend-name}} -t personal -d list -a {{page}}
```

2. Remote info

```bash
### Get info of the file from friend (public) ###
$ hls-filesharer-cli -f {{friend-name}} -d info -a {{filename}}
```

```bash
### Get info of the file from friend (personal) ###
$ hls-filesharer-cli -f {{friend-name}} -t personal -d info -a {{filename}}
```

3. Remote download

```bash
### Download file from friend (public) ###
$ hls-filesharer-cli -f {{friend-name}} -d download -a {{filename}}
```

```bash
### Download file from friend (personal) ###
$ hls-filesharer-cli -f {{friend-name}} -t personal -d download -a {{filename}}
```

4. Local list

```bash
### Get list of files (public) ###
$ hls-filesharer-cli -t local -d list -a {{page}}
```

```bash
### Get list of files from storage for friend (personal) ###
$ hls-filesharer-cli -f {{friend-name}} -t local -d list -a {{page}}
```

5. Local info

```bash
### Get info of the file (public) ###
$ hls-filesharer-cli -t local -d info -a {{filename}}
```

```bash
### Get info of the file from storage for friend (personal) ###
$ hls-filesharer-cli -f {{friend-name}} -t local -d info -a {{filename}}
```

6. Local download

```bash
### Download file (public) ###
$ hls-filesharer-cli -t local -d download -a {{filename}}
```

```bash
### Download file from storage for friend (personal) ###
$ hls-filesharer-cli -f {{friend-name}} -t local -d download -a {{filename}}
```

7. Local upload

```bash
### Upload file (public) ###
$ hls-filesharer-cli -t local -d upload -a {{path-to-file}}
```

```bash
### Upload file into friend's storage (personal) ###
$ hls-filesharer-cli -f {{friend-name}} -t local -d upload -a {{path-to-file}}
```

8. Local delete

```bash
### Delete file (public) ###
$ hls-filesharer-cli t local -d delete -a {{filename}}
```

```bash
### Delete file from friend's storage (personal) ###
$ hls-filesharer-cli -f {{friend-name}} -t local -d delete -a {{filename}}
```
