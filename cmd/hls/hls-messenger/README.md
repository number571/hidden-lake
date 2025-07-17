# HLS=messenger

> Hidden Lake Service (Messenger)

<img src="images/hls_messenger_logo.png" alt="hls_messenger_logo.png"/>

The `Hidden Lake Service (Messenger)` is a messenger based on the core of an anonymous network with theoretically provable anonymity of HLK. A feature of this messenger is the provision of anonymity of the fact of transactions (sending, receiving).

HLS=messenger is an application that implements a graphical user interface (GUI) on a browser-based HTML/CSS/JS display. Most of the code is based on the bootstrap library https://getbootstrap.com/. GUI is adapted to the size of the window, so it can be used both in a desktop and in a smartphone.

> More information about HLS=messenger in the [habr.com/ru/post/701488](https://habr.com/ru/post/701488/ "Habr HLS=messenger")

## Installation

```bash
$ go install github.com/number571/hidden-lake/cmd/hls/hls-messenger@latest
```

## How it works

Most of the code is a call to API functions from the HLK kernel. Thanks to this approach, implicit authorization of users is formed from the state of the anonymizing service.

<p align="center"><img src="images/hls_messenger_chat.gif" alt="hls_messenger_chat.gif"/></p>
<p align="center">Figure 1. Example of chat room in HLS=messenger.</p>

## Supported platforms

- Windows (x86_64, arm64)
- Linux (x86_64, arm64)
- MacOS (x86_64, arm64)

## Build and run

Default build and run

```bash 
$ go run ./cmd/hls/hls-messenger

> [INFO] 2023/06/03 15:30:31 HLS=messenger is running...
> ...
```

Open ports `9591` (HTTP, interface) and `9592` (HTTP, incoming).
Creates [`./hls-messenger.yml`](./hls-messenger.yml) and `./hls-messenger.db` files.

## Running options

```bash
$ hls-messenger --path /root
# path = path to config and database files
```

## Example

The example will involve (as well as in HLK) five nodes `node1_hlm, node2_hlm` and `middle_hla_tcp_1, middle_hla_tcp_2, middle_hla_tcp_3`. The three `HLA=tcp` nodes are only needed for communication between `node1_hlm` and `node2_hlm` nodes. Each of the remaining ones is a combination of HLK and HLS=messenger, where HLS=messenger plays the role of an application and services (as it was depicted in `Figure 3` HLK readme).

Build and run nodes
```bash
$ cd examples/messenger/routing
$ make
```

The output of the `middle_hls` node is similar to `Figure 4` (HLK).
Than open browser on `localhost:8080`. It is a `node1_hlm`. This node is a Bob.

<p align="center"><img src="images/hls_messenger_about.png" alt="hls_messenger_about.png"/></p>
<p align="center">Figure 2. Home page of the HLS=messenger application.</p>

To see the success of sending and receiving messages, you need to do all the same operations, but with `localhost:7070` as `node2_hlm`. This node will be Alice.

<p align="center"><img src="images/hls_messenger_logger.png" alt="hls_messenger_logger.png"/></p>
<p align="center">Figure 3. Log of the three nodes with request/response actions.</p>

> More example images about HLS=messenger pages in the [cmd/hls/hls-messenger/images](images "Path to HLS=messenger images")

## Pages

### About page

Base information about projects HLS=messenger and HLK with links to source.

<img src="images/v2/about.png" alt="about.png"/>

### Settings page

Information about public key and connections. Connections can be appended and deleted.

<img src="images/v2/settings.png" alt="settings.png"/>

### Friends page

Information about friends. Friends can be appended and deleted.

<img src="images/v2/friends.png" alt="friends.png"/>

### Chat page

Chat with friend. The chat is based on web sockets, so it can update messages in real time. Messages can be sent.

<img src="images/v2/chat.png" alt="chat.png"/>
