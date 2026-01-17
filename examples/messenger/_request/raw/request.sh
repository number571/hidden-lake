#!/bin/bash

## node2[localhost:8080](Alice) -> node1[localhost:7070](Bob)

PUSH_FORMAT='{
    "method":"POST",
    "host":"hls-messenger",
    "path":"/push",
    "body":"aGVsbG8sIHdvcmxkIQ=="
}';

curl -i -X PUT "http://localhost:8572/api/network/request?friend=Bob" --data "${PUSH_FORMAT}";
echo 
