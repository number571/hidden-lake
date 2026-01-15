#!/bin/bash

PUSH_FORMAT='{
    "method":"GET",
    "host":"hls-filesharer",
    "path":"/list?page=0"
}';

curl -i -X POST "http://localhost:8572/api/network/request?friend=Bob" --data "${PUSH_FORMAT}";
echo 
