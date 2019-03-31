# HiddenLake

## Characteristics
1. Based on p2p network.
2. Implemented end-to-end encryption.
3. Symmetric algorithm = AES256-CBC.
4. Asymmetric algorithm = RSA2048-OAEP.

## Compile and run client:
```
$ go build -o main main.go
$ ./main --address 127.0.0.1:8080 --interface
```
![Image alt](https://github.com/Number571/HiddenLake/raw/master/images/HiddenLake_CLI_1.png)

## GUI work in browser on port :7545:
```
$ firefox --new-window 127.0.0.1:7545
```
![Image alt](https://github.com/Number571/HiddenLake/raw/master/images/HiddenLake_GUI_1.png)
