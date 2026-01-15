#!/bin/bash

## node2[localhost:7070](Bob) -> node1[localhost:8080](Alice)

PUSH_FORMAT='{
    "method":"POST",
    "host":"hls-messenger",
    "path":"/push",
    "body":"aGVsbG8sIHdvcmxkIQ=="
}';

curl -i -X PUT "http://localhost:7572/api/network/request?friend=Alice" --data "${PUSH_FORMAT}";
echo 
