#!/bin/bash

go run ./../../cmd/hls/hls-filesharer/client/hls-filesharer-cli -s localhost:8541 -t public -f Bob -d info -a image.jpg
