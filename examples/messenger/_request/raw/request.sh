#!/bin/bash

## node2[localhost:7070](Bob) -> node1[localhost:8080](Alice)

PUSH_FORMAT='{
    "receiver":"Alice",
    "req_data":{
        "method":"POST",
        "host":"hls-messenger",
        "path":"/push",
        "body":"aGVsbG8sIHdvcmxkIQ=="
    }
}';

curl -i -X PUT http://localhost:7572/api/network/request --data "${PUSH_FORMAT}";
echo 
