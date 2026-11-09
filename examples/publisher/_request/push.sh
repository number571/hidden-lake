#!/bin/bash

echo "<Alice's side>"
go run ./../../cmd/hls/hls-publisher/client/hls-publisher-cli -s localhost:8080 -f Bob -d push -a "$1"
