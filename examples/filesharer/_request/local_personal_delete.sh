#!/bin/bash

go run ./../../cmd/hls/hls-filesharer/client/hls-filesharer-cli -s localhost:7541 -t local -f Alice -d delete -a image.jpg
