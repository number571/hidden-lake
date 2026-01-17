# HLS=filesharer

> Hidden Lake Service (Filesharer)

<img src="images/hls_filesharer_logo.png" alt="hls_filesharer_logo.png"/>

The `Hidden Lake Service (Filesharer)` is a file sharing service based on the anonymous network core (HLK) with theoretically provable anonymity. A feature of this file sharing service is the anonymity of the fact of transactions (file downloads), taking into account the existence of a global observer.

HLS=filesharer is an application that implements a graphical user interface (GUI) on a browser-based HTML/CSS/JS display. Most of the code is based on the bootstrap library https://getbootstrap.com/. GUI is adapted to the size of the window, so it can be used both in a desktop and in a smartphone.

> More information about HLS=filesharer in the [habr.com/ru/articles/789968](https://habr.com/ru/articles/789968/ "Habr HLS=filesharer")

## Installation

```bash
$ go install github.com/number571/hidden-lake/cmd/hls/hls-filesharer@latest
```

## How it works

Most of the code is a call to API functions from the HLK kernel. Thanks to this approach, implicit authorization of users is formed from the state of the anonymizing service.

Unlike applications such as HLK and HLS=messenger, the HLS=filesharer application does not have a database. Instead, the storage is used, represented by the usual `hls-filesharer.stg` directory.

<p align="center"><img src="images/hls_filesharer_download.gif" alt="hls_filesharer_download.gif"/></p>
<p align="center">Figure 1. Example of download file in HLS=filesharer (x2 speed).</p>

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
1. GET  /api/index                  | params = []
                                    |> description = get name of service
2. GET  /api/storage/list           | params = ["friend":string,"page":uint64]
                                    |> description = get list of files from storage
3. GET  /api/storage/file/info      | params = ["friend":string,"name":string]
                                    |> description = get info of the file by name
4. GET  /api/storage/file/download  | params = ["friend":string,"name":string]
                                    |> description = download file content by name
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

### 2. /api/storage/list

#### 2.1. GET Request

```bash
curl -i -X GET "http://localhost:9541/api/storage/list?friend=Bob&page=0"
```

#### 2.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 16 Jan 2026 21:25:00 GMT
Content-Length: 280

[{"name":"example.txt","hash":"7d0c64e050a2c31cd2d5266b2923ca51b95e97e2dedfc39e4ce220b477683975ba032c6c3141bad8442af4943f91ac43","size":14},{"name":"image.jpg","hash":"7bfd88d546b47b60dba2cd5f5ff8f2ccc19daf348640d5f5d3381bc3f54306f1c7984679c235b4c4485b56eec2ad4977","size":17792}]
```

### 3. /api/storage/file/info

#### 3.1. GET Request

```bash
curl -i -X GET "http://localhost:9541/api/storage/file/info?friend=Bob&name=example.txt"
```

#### 3.1. GET Response

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 16 Jan 2026 21:25:55 GMT
Content-Length: 138

{"name":"example.txt","hash":"7d0c64e050a2c31cd2d5266b2923ca51b95e97e2dedfc39e4ce220b477683975ba032c6c3141bad8442af4943f91ac43","size":14}
```

### 4. /api/storage/file/download

#### 4.1. GET Request

```bash
curl -i -X GET "http://localhost:9541/api/storage/file/download?friend=Bob&name=example.txt"
```

#### 4.1. GET Response

```
HTTP/1.1 200 OK
Date: Fri, 16 Jan 2026 21:26:55 GMT
Content-Length: 14
Content-Type: text/plain; charset=utf-8

hello, world!
```
