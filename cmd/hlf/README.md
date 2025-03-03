# HLF

> Hidden Lake Filesharer

<img src="images/hlf_logo.png" alt="hlf_logo.png"/>

The `Hidden Lake Filesharer` is a file sharing service based on the anonymous network core (HLS) with theoretically provable anonymity. A feature of this file sharing service is the anonymity of the fact of transactions (file downloads), taking into account the existence of a global observer.

HLF is an application that implements a graphical user interface (GUI) on a browser-based HTML/CSS/JS display. Most of the code is based on the bootstrap library https://getbootstrap.com/. GUI is adapted to the size of the window, so it can be used both in a desktop and in a smartphone.

> More information about HLF in the [habr.com/ru/articles/789968](https://habr.com/ru/articles/789968/ "Habr HLF")

## Installation

```bash
$ go install github.com/number571/hidden-lake/cmd/hlf@latest
```

## How it works

Most of the code is a call to API functions from the HLS kernel. Thanks to this approach, implicit authorization of users is formed from the state of the anonymizing service.

Unlike applications such as HLS and HLM, the HLF application does not have a database. Instead, the storage is used, represented by the usual `hlf.stg` directory.

<p align="center"><img src="images/hlf_download.gif" alt="hlf_download.gif"/></p>
<p align="center">Figure 1. Example of download file in HLF (x2 speed).</p>

File transfer is limited by the bandwidth of HLS itself. If we take into account that the packet generation period is `5 seconds`, then it will take about 10 seconds to complete the request-response cycle. HLS also limits the size of transmitted packets. If we assume that the limit is `8KiB`, taking into account the existing ~4KiB headers, then the transfer rate is defined as `4KiB/10s` or `410B/1s`.

## Supported platforms

- Windows (x86_64, arm64)
- Linux (x86_64, arm64)
- MacOS (x86_64, arm64)

## Build and run

Default build and run

```bash 
$ go run ./cmd/hlf

> [INFO] 2023/06/03 15:30:31 HLF is running...
> ...
```

Open ports `9541` (HTTP, interface) and `9542` (HTTP, incoming).
Creates [`./hlf.yml`](./hlf.yml) and `./hlf.stg` files.
The directory `hlf.stg` stores all shared/loaded files. 

## Running options

```bash
$ hlf --path /root
# path = path to config and storage
```

## Example

The example will involve two nodes `node1_hlf, node2_hlf` and three repeaters `middle_hla_tcp_1, middle_hla_tcp_2, middle_hla_tcp_3`. Both nodes are a combination of HLS and HLF, where HLF plays the role of an application and services (as shown in `Figure 3` of the HLS readme). The three remaining nodes are used only for the successful connection of the two main nodes. In other words, `HLA=tcp` nodes are traffic relay nodes.

Build and run nodes
```bash
$ cd examples/filesharer/routing
$ make
```

Than open browser on `localhost:8080`. It is a `node1_hlf`. This node is a Alice.

<p align="center"><img src="images/hlf_about.png" alt="hlf_about.png"/></p>
<p align="center">Figure 2. Home page of the HLF application.</p>

To see the another side of communication, you need to do all the same operations, but with `localhost:7070` as `node2_hlf`. This node will be Bob.

> More example images about HLF pages in the [cmd/hlf/images](images "Path to HLF images")
