# HiddenLake

## Characteristics:
1. Hybrid network.
2. Implemented onion routing.
3. Implemented P2P/F2F connection.
4. Implemented end-to-end encryption.
5. Implemented client-server connection.
6. Symmetric algorithm = AES256-CBC.
7. Asymmetric algorithm = RSA2048-OAEP.
8. Cryptographic hash functions = MD5/SHA256.

### // Home page
<img src="/images/HiddenLake_GUI_1.png" alt="GUI_1"/>

## Components are used:
1. Go: go-sqlite3
2. JS: jquery

Go version should be >= 1.10

## Component installation:
```
$ go get github.com/mattn/go-sqlite3
```

## Compile:
```
$ cd HiddenLake/
$ go build -ldflags "-w -s" main.go
```

## Commands in the start client with parameters:
1. [--login, -l] = set login (first run is signup)
2. [--password, -p] = set password (first run is signup)
3. [--address, -a] = set address ipv4:port

## Commands in the start client without parameters:
1. [--interface, -i] = run GUI interface in browser on port 7545
2. [--help, -h] = get information about client
3. [--f2f, -f] = run F2F connection
4. [--delete, -d] = delete database and archive files with multiple overwriting
4. [--delete-database, -dd] = delete database file with multiple overwriting
4. [--delete-archive, -da] = delete archive files with multiple overwriting

## Run CLI client:
```
$ ./main --login "user" --password "hello, world" --address 127.0.0.1:8080
```

### // CLI
<img src="/images/HiddenLake_CLI_1.png" alt="CLI_1"/>

## Commands in CLI client for all users:
1. [:exit] = exit from client
2. [:help] = get information about client
3. [:mode] = on/off F2F connection
4. [:interface] = on/off GUI interface

## Commands in CLI client if not authorized:
1. [:login] = set login (first run is signup)
2. [:password] = set password (first run in signup)
3. [:enter] = authorization from the entered login and password
4. [:address] = set address ipv4:port

## Commands in CLI client if authorized:
1.  [:whoami] = get hashname
2.  [:logout] = logout from authorized user
3.  [:network] = get list of connections
4.  [:send] = send local message to another user
5.  [:email] = read or write email to another user
6.  [:archive] = get list or download files from archive another user
7.  [:history] = get local/global messages or delete messages
8.  [:connect] = connect to another user
9.  [:disconnect] = disconnect from user
10. [] = send global message to another users

## Run CLI/GUI client:
> GUI work in browser on port 7545

```
$ ./main --interface
$ firefox --new-window 127.0.0.1:7545
```
### // Login page
<img src="/images/HiddenLake_GUI_9.png" alt="GUI_9"/>

## [HiddenLake]
