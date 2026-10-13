#!/bin/bash

ROOT_DIR=./../..

# upload files
go run $ROOT_DIR/cmd/hls/hls-filesharer/client/hls-filesharer-cli -s localhost:7541 -t local -d upload -a $ROOT_DIR/examples/filesharer/_storage/image.jpg
go run $ROOT_DIR/cmd/hls/hls-filesharer/client/hls-filesharer-cli -s localhost:7541 -t local -d upload -a $ROOT_DIR/examples/filesharer/_storage/example.txt

PUSH_FORMAT='{
    "method":"GET",
    "host":"hls-filesharer",
    "path":"/list?page=0"
}';

curl -i -X POST "http://localhost:8572/api/network/request?friend=Bob" --data "${PUSH_FORMAT}";
echo 
