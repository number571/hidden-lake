# HLS=filesharer

> Hidden Lake Service (Filesharer)

<img src="images/hls_filesharer_logo.png" alt="hls_filesharer_logo.png"/>

The `Hidden Lake Service (Filesharer)` is a file sharing service based on the anonymous network core (HLK) with theoretically provable anonymity. A feature of this file sharing service is the anonymity of the fact of transactions (file downloads), taking into account the existence of a global observer.

> More information about HLS=filesharer in the [habr.com/ru/articles/789968](https://habr.com/ru/articles/789968/ "Habr HLS=filesharer")

## Installation

```bash
$ go install github.com/number571/hidden-lake/cmd/hls/hls-filesharer@latest
```

## How it works

Most of the code is a call to API functions from the HLK kernel. Thanks to this approach, implicit authorization of users is formed from the state of the anonymizing service. Unlike applications such as HLK and HLS=notifier, the HLS=filesharer application does not have a database. Instead, the storage is used, represented by the usual `hls-filesharer.stg` directory.

File transfer is limited by the bandwidth of HLK itself. If we take into account that the packet generation period is `5 seconds`, then it will take about 10 seconds to complete the request-response cycle. HLK also limits the size of transmitted packets. If we assume that the limit is `8KiB`, taking into account the existing ~4KiB headers, then the transfer rate is defined as `4KiB/10s` or `410B/1s`.

## Supported platforms

- Windows (x86_64, arm64)
- Linux (x86_64, arm64)
- MacOS (x86_64, arm64)

## Build and run

Default build and run

```bash 
$ go run ./cmd/hls/hls-filesharer

> [INFO] 2023/06/03 15:30:31 HLS=filesharer is running...
> ...
```

Open ports `9541` (HTTP, internal) and `9542` (HTTP, incoming).
Creates [`./hls-filesharer.yml`](./hls-filesharer.yml) and `./hls-filesharer.stg` files.
The directory `hls-filesharer.stg` stores all shared/loaded files. 

## Running options

```bash
$ hls-filesharer --path /root
# path = path to config and storage
```

## Example

The example will involve two nodes `node1_hlf, node2_hlf` and three repeaters `middle_hla_tcp_1, middle_hla_tcp_2, middle_hla_tcp_3`. Both nodes are a combination of HLK and HLS=filesharer, where HLS=filesharer plays the role of an application and services (as shown in `Figure 3` of the HLK readme). The three remaining nodes are used only for the successful connection of the two main nodes. In other words, `HLA=tcp` nodes are traffic relay nodes.

Build and run nodes
```bash
$ cd examples/filesharer/routing
$ make
```

Than run command
```bash
$ cd examples/filesharer
$ ./_request/raw/request.sh
```

Got response
```json
{"code":200,"body":"WyJleGFtcGxlLnR4dCIsImltYWdlLmpwZyJd"}
```

## HLS/HLK API

### 1. /info

- `200` - return file info
- `400` - invalid input params (name, personal)
- `404` - file not found
- `405` - invalid method (not GET)
- `500` - load file / info
- `502` - get aliasName HLK

```json
// header request size = (4 + 76) = 80B
// body request size (with query params) can be = (144 - 80) = 64B
[uint32]{"method":"GET","host":"hls-filesharer","path":"/info?name=&personal=false"}[body]

// header response size = (4 + 12) = 16B
// body response size (with name and size) can be = (144 - 16 - 70) = 58B
// INFORMATION: hash = trim(sha384) => size(20B)
// RECOMMENDATION: file_name_size=20B
[uint32]{"code":200}[body]
[body]={"name":"","size":0,"hash":"0000000000000000000000000000000000000000"}
```

### 2. /list

- `200` - return list of file names
- `400` - invalid input params (page, personal)
- `405` - invalid method (not GET)
- `500` - load list of files
- `502` - get aliasName HLK

```json
// header request size = (4 + 76) = 80B
// body request size (with query params) can be = (144 - 80) = 64B
[uint32]{"method":"GET","host":"hls-filesharer","path":"/list?page=&personal=false"}[body]

// header response size = (4 + 12) = 16B
// body response size (with list of JSON strings) can be = (144 - 16 - 2) = 126B
// RECOMMENDATION: page_offset=5 and file_name_size<=20B
[uint32]{"code":200}[body]
[body]=[]
```

### 3. /load

- `200` - return chunk of file
- `400` - invalid input params (name, chunk, personal)
- `404` - file not found
- `405` - invalid method (not GET)
- `411` - invalid file chunk number
- `500` - open file, load file chunk
- `502` - get aliasName HLK, get limit size on load from HLK

```json
// header request size = (4 + 83) = 87B
// body request size (with query params) can be = (144 - 87) = 57B
[uint32]{"method":"GET","host":"hls-filesharer","path":"/load?name=&chunk=&personal=false"}[body]

// header response size = (4 + 46) = 50B
// body response size can be = (144 - 50) = 94B
// INFORMATION: hash = trim(sha384) => size(8B)
[uint32]{"code":200,"head":{"F-H":"0000000000000000"}}[body]
```

## HLS API

```
1. GET                  /api/index              | params = []
                                                |> description = get name of service
2. GET                  /api/remote/list        | params = ["friend":string,"page":uint64,"personal":?bool]
                                                |> description = get list of files from storage
3. GET                  /api/remote/list/proc   | params = []
                                                |> description = get list of file download processes
4. GET, DELETE          /api/remote/file        | params = ["friend":string,"name":string,"personal":?bool]
                                                |> description = download / delete file content by name
5. GET, DELETE          /api/remote/file/proc   | params = ["friend":string,"name":string,"personal":?bool]
                                                |> description = view / cancel file download process
6. GET                  /api/remote/file/info   | params = ["friend":string,"name":string,"personal":?bool]
                                                |> description = get info of the file by name
7. GET                  /api/local/list         | params = ["friend":?string,"page":uint64]
                                                |> description = get list of files from storage
8. GET, POST, DELETE    /api/local/file         | params = ["friend":?string,"name":string]
                                                |> description = download / upload / delete file content by name
9. GET                  /api/local/file/info    | params = ["friend":?string,"name":string]
                                                |> description = get info of the file by name
```

### 1. /api/index

- `200` - return "hidden-lake-service=filesharer"

#### 1.1. GET Request

```bash
curl -i -X GET http://localhost:9541/api/index
```

#### 1.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Fri, 16 Jan 2026 21:24:23 GMT
Content-Length: 30

hidden-lake-service=filesharer
```

### 2. /api/remote/list

- `200` - return list of files
- `400` - invalid input params (friend, page, personal)
- `405` - invalid method (not GET)
- `418` - invalid response from node
- `502` - fetch request HLK

#### 2.1. GET Request

```bash
curl -i -X GET "http://localhost:9541/api/remote/list?friend=Bob&page=0&personal=false"
```

#### 2.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 16 Jan 2026 21:25:00 GMT
Content-Length: 27

["example.txt","image.jpg"]
```

### 3. /api/remote/list/proc

- `200` - download process info list
- `405` - invalid method (not GET)

#### 3.1. GET Request

```bash
curl -i -X GET "http://localhost:9541/api/remote/list/proc"
```

#### 3.1. GET Response

```
HTTP/1.1 200 OK
F-H: 7d0c64e050a2c31c
Date: Mon, 19 Jan 2026 08:05:06 GMT
Content-Length: 101
Content-Type: application/json; charset=utf-8

[{{"friend":"Bob","name":"example.txt","personal":"false"},{"incindex":0,"download":0,"filesize":0}}]
```

### 4. /api/remote/file

- `200` - return file content / delete file
- `202` - download process already exist 
- `400` - invalid input params (friend, name, personal)
- `405` - invalid method (not GET, DELETE)
- `418` - invalid response from node
- `500` - mkdir storage, remove file, build/read stream, open file
- `502` - fetch request HLK

#### 4.1. GET Request

```bash
curl -i -X GET "http://localhost:9541/api/remote/file?friend=Bob&name=example.txt&personal=false"
```

#### 4.1. GET Response

```
HTTP/1.1 200 OK
Hls-Filesharer-In-Process: -
Date: Mon, 19 Jan 2026 08:05:06 GMT
Content-Length: 14
Content-Type: text/plain; charset=utf-8

hello, world!
```

#### 4.2. DELETE Request

```bash
curl -i -X DELETE "http://localhost:9541/api/remote/file?friend=Bob&name=example.txt&personal=false"
```

#### 4.2. DELETE Response

```
HTTP/1.1 200 OK
Date: Mon, 19 Jan 2026 08:05:06 GMT
Content-Length: 20
Content-Type: text/plain; charset=utf-8

success: delete file
```

### 5. /api/remote/file/proc

- `200` - download process info / cancel download
- `400` - invalid input params (friend, name, personal)
- `404` - process not found
- `405` - invalid method (not GET, DELETE)

#### 5.1. GET Request

```bash
curl -i -X GET "http://localhost:9541/api/remote/file/proc?friend=&name=example.txt&personal=false"
```

#### 5.1. GET Response

```
HTTP/1.1 200 OK
Date: Sun, 01 Feb 2026 16:11:23 GMT
Content-Length: 99
Content-Type: application/json; charset=utf-8

{{"friend":"Bob","name":"example.txt","personal":"false"},{"incindex":0,"download":0,"filesize":0}}
```

#### 5.2. DELETE Request

```bash
curl -i -X DELETE "http://localhost:9541/api/remote/file/proc?friend=Alice&name=example.txt&personal=false"
```

#### 5.2. DELETE Response

```
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Sun, 01 Feb 2026 16:15:34 GMT
Content-Length: 0

success: cancel download
```

### 6. /api/remote/file/info

- `200` - return file info
- `400` - invalid input params (friend, name, personal)
- `405` - invalid method (not GET)
- `418` - invalid response from node
- `502` - fetch request HLK

#### 6.1. GET Request

```bash
curl -i -X GET "http://localhost:9541/api/remote/file/info?friend=Bob&name=example.txt&personal=false"
```

#### 6.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 16 Jan 2026 21:25:55 GMT
Content-Length: 82

{"name":"example.txt","size":14,"hash":"7d0c64e050a2c31cd2d5266b2923ca51b95e97e2"}
```

### 7. /api/local/list

- `200` - return list of files
- `400` - invalid input params (friend, page)
- `405` - invalid method (not GET)
- `500` - get list of files

#### 7.1. GET Request

```bash
curl -i -X GET "http://localhost:9541/api/local/list?friend=&page=0"
```

#### 7.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sun, 01 Feb 2026 16:09:56 GMT
Content-Length: 27

["example.txt","image.jpg"]
```

### 8. /api/local/file

- `200` - return file content / delete file / upload file
- `400` - invalid input params (friend, name)
- `404` - file not found
- `405` - invalid method (not GET, POST, DELETE)
- `500` - mkdir, create/open/remove file

#### 8.1. GET Request

```bash
curl -i -X GET "http://localhost:9541/api/local/file?friend=&name=example.txt"
```

#### 8.1. GET Response

```
HTTP/1.1 200 OK
Date: Sun, 01 Feb 2026 16:11:23 GMT
Content-Length: 14
Content-Type: text/plain; charset=utf-8

hello, world!
```

#### 8.2. POST Request

```bash
curl -i -X POST "http://localhost:9541/api/local/file?friend=Alice&name=example.txt" --data 'hello, world!'
```

#### 8.2. POST Response

```
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Sun, 01 Feb 2026 16:12:18 GMT
Content-Length: 20

success: upload file
```

#### 8.3. DELETE Request

```bash
curl -i -X DELETE "http://localhost:9541/api/local/file?friend=Alice&name=example.txt"
```

#### 8.3. DELETE Response

```
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Sun, 01 Feb 2026 16:15:34 GMT
Content-Length: 0

success: delete file
```

### 9. /api/local/file/info

- `200` - return file info
- `400` - invalid input params (friend, name)
- `404` - file not found
- `405` - invalid method (not GET)
- `500` - open/read file

#### 9.1. GET Request

```bash
curl -i -X GET "http://localhost:9541/api/local/file/info?friend=&name=example.txt"
```

#### 9.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sun, 01 Feb 2026 16:17:12 GMT
Content-Length: 138

{"name":"example.txt","size":14,"hash":"7d0c64e050a2c31cd2d5266b2923ca51b95e97e2"}
```
