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

Most of the code is a call to API functions from the HLK kernel. Thanks to this approach, implicit authorization of users is formed from the state of the anonymizing service. Unlike applications such as HLK and HLS=messenger, the HLS=filesharer application does not have a database. Instead, the storage is used, represented by the usual `hls-filesharer.stg` directory.

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
$ cd examples/messenger
$ ./_request/raw/request.sh
```

Got response
```json
{"code":200,"head":{"Content-Type":"application/json"},"body":"W3sibmFtZSI6ImV4YW1wbGUudHh0IiwiaGFzaCI6IjdkMGM2NGUwNTBhMmMzMWNkMmQ1MjY2YjI5MjNjYTUxYjk1ZTk3ZTJkZWRmYzM5ZTRjZTIyMGI0Nzc2ODM5NzViYTAzMmM2YzMxNDFiYWQ4NDQyYWY0OTQzZjkxYWM0MyIsInNpemUiOjE0fSx7Im5hbWUiOiJpbWFnZS5qcGciLCJoYXNoIjoiN2JmZDg4ZDU0NmI0N2I2MGRiYTJjZDVmNWZmOGYyY2NjMTlkYWYzNDg2NDBkNWY1ZDMzODFiYzNmNTQzMDZmMWM3OTg0Njc5YzIzNWI0YzQ0ODViNTZlZWMyYWQ0OTc3Iiwic2l6ZSI6MTc3OTJ9XQ=="}
```

## HLS API

```
1. GET                  /api/index              | params = []
                                                |> description = get name of service
2. GET                  /api/remote/list        | params = ["friend":string,"page":uint64,"personal":?bool]
                                                |> description = get list of files from storage
3. GET, DELETE          /api/remote/file        | params = ["friend":string,"name":string,"personal":?bool]
                                                |> description = download file content by name
4. GET                  /api/remote/file/info   | params = ["friend":string,"name":string,"personal":?bool]
                                                |> description = get info of the file by name
5. GET                  /api/local/list         | params = ["friend":?string,"page":uint64]
                                                |> description = get list of files from storage
6. GET, POST, DELETE    /api/local/file         | params = ["friend":?string,"name":string]
                                                |> description = upload / delete file content by name
7. GET                  /api/local/file/info    | params = ["friend":?string,"name":string]
                                                |> description = get info of the file by name
```

### 1. /api/index

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

#### 2.1. GET Request

```bash
curl -i -X GET "http://localhost:9541/api/remote/list?friend=Bob&page=0&personal=false"
```

#### 2.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 16 Jan 2026 21:25:00 GMT
Content-Length: 280

[{"name":"example.txt","size":14,"hash":"7d0c64e050a2c31cd2d5266b2923ca51b95e97e2dedfc39e4ce220b477683975ba032c6c3141bad8442af4943f91ac43"},{"name":"image.jpg","size":17792,"hash":"7bfd88d546b47b60dba2cd5f5ff8f2ccc19daf348640d5f5d3381bc3f54306f1c7984679c235b4c4485b56eec2ad4977"}]
```

### 3. /api/remote/file

#### 3.1. GET Request

```bash
curl -i -X GET "http://localhost:9541/api/remote/file?friend=Bob&name=example.txt&personal=false"
```

#### 3.1. GET Response

```
HTTP/1.1 200 OK
Hls-Filesharer-File-Hash: 7d0c64e050a2c31cd2d5266b2923ca51b95e97e2dedfc39e4ce220b477683975ba032c6c3141bad8442af4943f91ac43
Date: Mon, 19 Jan 2026 08:05:06 GMT
Content-Length: 14
Content-Type: text/plain; charset=utf-8

hello, world!
```

#### 3.2. DELETE Request

```bash
curl -i -X DELETE "http://localhost:9541/api/remote/file?friend=Bob&name=example.txt&personal=false"
```

#### 3.2. DELETE Response

```
HTTP/1.1 200 OK
Date: Mon, 19 Jan 2026 08:05:06 GMT
Content-Length: 20
Content-Type: text/plain; charset=utf-8

success: delete file
```

### 4. /api/remote/file/info

#### 4.1. GET Request

```bash
curl -i -X GET "http://localhost:9541/api/remote/file/info?friend=Bob&name=example.txt&personal=false"
```

#### 4.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 16 Jan 2026 21:25:55 GMT
Content-Length: 138

{"name":"example.txt","size":14,"hash":"7d0c64e050a2c31cd2d5266b2923ca51b95e97e2dedfc39e4ce220b477683975ba032c6c3141bad8442af4943f91ac43"}
```

### 5. /api/local/list

#### 5.1. GET Request

```bash
curl -i -X GET "http://localhost:9541/api/local/list?friend=&page=0"
```

#### 5.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Sun, 01 Feb 2026 16:09:56 GMT
Content-Length: 280

[{"name":"example.txt","size":14,"hash":"7d0c64e050a2c31cd2d5266b2923ca51b95e97e2dedfc39e4ce220b477683975ba032c6c3141bad8442af4943f91ac43"},{"name":"image.jpg","size":17792,"hash":"7bfd88d546b47b60dba2cd5f5ff8f2ccc19daf348640d5f5d3381bc3f54306f1c7984679c235b4c4485b56eec2ad4977"}]
```

### 6. /api/local/file

#### 6.1. GET Request

```bash
curl -i -X GET "http://localhost:9541/api/local/file?friend=&name=example.txt"
```

#### 6.1. GET Response

```
HTTP/1.1 200 OK
Date: Sun, 01 Feb 2026 16:11:23 GMT
Content-Length: 14
Content-Type: text/plain; charset=utf-8

hello, world!
```

#### 6.2. POST Request

```bash
curl -i -X POST "http://localhost:9541/api/local/file?friend=Alice&name=example.txt" --data 'hello, world!'
```

#### 6.2. POST Response

```
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Sun, 01 Feb 2026 16:12:18 GMT
Content-Length: 20

success: upload file
```

#### 6.3. DELETE Request

```bash
curl -i -X DELETE "http://localhost:9541/api/local/file?friend=Alice&name=example.txt"
```

#### 6.3. DELETE Response

```
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Sun, 01 Feb 2026 16:15:34 GMT
Content-Length: 0

success: delete file
```

### 7. /api/local/file/info

#### 7.1. GET Request

```bash
curl -i -X GET "http://localhost:9541/api/local/file/info?friend=&name=example.txt"
```

#### 7.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: text/plain
Date: Sun, 01 Feb 2026 16:17:12 GMT
Content-Length: 138

{"name":"example.txt","size":14,"hash":"7d0c64e050a2c31cd2d5266b2923ca51b95e97e2dedfc39e4ce220b477683975ba032c6c3141bad8442af4943f91ac43"}
```
