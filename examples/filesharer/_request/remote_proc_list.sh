#!/bin/bash

echo "<Alice's side>"
go run ./../../cmd/hls/hls-filesharer/client/hls-filesharer-cli -s localhost:8541 -f Bob -d proc
