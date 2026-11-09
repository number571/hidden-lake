#!/bin/bash

echo "<Bob's side>"
go run ./../../cmd/hls/hls-publisher/client/hls-publisher-cli -s localhost:7070 -f Alice -r -d load-by-uid -a "$1"
