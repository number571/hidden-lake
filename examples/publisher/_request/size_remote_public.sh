#!/bin/bash

echo "<Bob's side>"
go run ./../../cmd/hls/hls-publisher/client/hls-publisher-cli -s localhost:7070 -r -p -f Alice -d size
