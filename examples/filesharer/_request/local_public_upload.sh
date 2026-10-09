#!/bin/bash

echo "<Bob's side>"
go run ./../../cmd/hls/hls-filesharer/client/hls-filesharer-cli -s localhost:7541 -t local -d upload -a ./../../examples/filesharer/_storage/image.jpg
go run ./../../cmd/hls/hls-filesharer/client/hls-filesharer-cli -s localhost:7541 -t local -d upload -a ./../../examples/filesharer/_storage/example.txt
